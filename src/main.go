package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ignoring missing .env file")
	}
}

type Issue struct {
	Publisher   int
	Publication int
	Issue       int
}

var buckets = map[string]func(Issue) string{
	//"media.magtab.com": func(issue Issue) string {
	//	return fmt.Sprintf("editoras/%d/titulos/%d/edicoes/%d", issue.Publisher, issue.Publication, issue.Issue)
	//},
	"revistas.magtab.com": func(issue Issue) string {
		return fmt.Sprintf("%d/%d/%d", issue.Publisher, issue.Publication, issue.Issue)
	},
}

func main() {
	awsConfig := getAWSConfig()
	ociConfig := getOCIConfig()
	issue := Issue{Publisher: 9, Publication: 17, Issue: 21650}
	copyIssueFiles(awsConfig, ociConfig, issue)
}

func getOCIConfig() aws.Config {
	accessKey := os.Getenv("DEST_ACCESS_KEY")
	secretKey := os.Getenv("DEST_SECRET_KEY")
	endpointURL := os.Getenv("DEST_ENDPOINT_URL")
	region := os.Getenv("DEST_REGION")

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           endpointURL,
			SigningRegion: region,
		}, nil
	})

	config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion(region),
		config.WithEndpointResolver(aws.EndpointResolver(customResolver)),
	)
	if err != nil {
		log.Fatalf("Failed to load OCI configuration, %v", err)
	}
	return config
}

func getAWSConfig() aws.Config {
	accessKey := os.Getenv("SOURCE_ACCESS_KEY")
	secretKey := os.Getenv("SOURCE_SECRET_KEY")
	region := os.Getenv("SOURCE_REGION")

	config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS configuration, %v", err)
	}
	return config
}

func copyIssueFiles(awsConfig aws.Config, ociConfig aws.Config, issue Issue) {
	for bucketName, prefixFunc := range buckets {
		prefix := prefixFunc(issue)
		ociItems := getObjectKeys(ociConfig, bucketName, prefix)
		fmt.Println("=", bucketName)
		fmt.Println("- OCI Items")
		copier1 := copyObjects(ociConfig, bucketName, ociItems)
		copier2 := copyObjects(ociConfig, bucketName, ociItems)
		copier3 := copyObjects(ociConfig, bucketName, ociItems)
		for x := range merge(copier1, copier2, copier3) {
			fmt.Println(x)
		}

		//awsItems := getObjectKeys(awsConfig, bucketName, prefix)
		//fmt.Println("- AWS Items")
		//for x := range awsItems {
		//	fmt.Println(x)
		//}
	}
}

func merge(channels ...<-chan os.File) <-chan os.File {
	var waitGroup sync.WaitGroup
	output := make(chan os.File)

	partialOutput := func(channel <-chan os.File) {
		for key := range channel {
			output <- key
		}
		waitGroup.Done()
	}
	waitGroup.Add(len(channels))
	for _, channel := range channels {
		go partialOutput(channel)
	}

	go func() {
		waitGroup.Wait()
		close(output)
	}()
	return output
}

func getObjectKeys(config aws.Config, bucketName string, prefix string) <-chan string {
	client := s3.NewFromConfig(config)
	params := &s3.ListObjectsV2Input{Bucket: &bucketName, Prefix: &prefix}
	paginator := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		o.Limit = 10
	})
	channel := make(chan string)
	go func() {
		var i int
		for paginator.HasMorePages() {
			i++

			page, err := paginator.NextPage(context.TODO())
			if err != nil {
				log.Fatalf("failed to get page %v, %v", i, err)
			}

			// Log the objects found
			for _, obj := range page.Contents {
				channel <- *obj.Key
			}

			if i > 0 {
				break
			}
		}
		close(channel)
	}()
	return channel
}

func copyObjects(config aws.Config, bucketName string, input <-chan string) <-chan os.File {
	client := s3.NewFromConfig(config)
	output := make(chan os.File)
	go func() {
		for key := range input {
			requestInput := s3.GetObjectInput{
				Bucket: &bucketName,
				Key:    &key,
			}
			pathToFile := os.TempDir() + "/" + path.Base(key)
			fmt.Println(pathToFile)
			f, _ := os.Create(pathToFile)
			downloader := manager.NewDownloader(client)
			_, err := downloader.Download(context.TODO(), f, &requestInput)
			if err != nil {
				log.Fatalf("failed to download %s", key)
			}
			output <- *f
		}
		close(output)
	}()
	return output

}
