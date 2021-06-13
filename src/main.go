package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

var hardLimit int32 = 6

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ignoring missing .env file")
	}
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
	accessKey := os.Getenv("OCI_ACCESS_KEY")
	secretKey := os.Getenv("OCI_SECRET_KEY")
	endpointURL := os.Getenv("OCI_ENDPOINT_URL")
	region := os.Getenv("OCI_REGION")

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
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	region := os.Getenv("AWS_REGION")

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
		copier1 := copyObjects(ociConfig, awsConfig, bucketName, ociItems)
		copier2 := copyObjects(ociConfig, awsConfig, bucketName, ociItems)
		copier3 := copyObjects(ociConfig, awsConfig, bucketName, ociItems)
		for item := range mergeCopiers(copier1, copier2, copier3) {
			fmt.Println("Input", item)
		}
	}
}

func mergeCopiers(channels ...<-chan string) <-chan string {
	var waitGroup sync.WaitGroup
	output := make(chan string)

	partialOutput := func(channel <-chan string) {
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
		o.Limit = hardLimit
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

func copyObjects(ociConfig aws.Config, awsConfig aws.Config, bucketName string, input <-chan string) <-chan string {
	ociClient := s3.NewFromConfig(ociConfig)
	awsClient := s3.NewFromConfig(awsConfig)
	output := make(chan string)
	go func() {
		for key := range input {
			requestInput := s3.GetObjectInput{
				Bucket: &bucketName,
				Key:    &key,
			}

			reader, writer := io.Pipe()

			downloader := manager.NewDownloader(ociClient)
			downloader.Concurrency = 1
			go func() {
				defer writer.Close()
				downloader.Download(context.TODO(), FakeWriterAt{writer}, &requestInput)
			}()

			uploadInput := s3.PutObjectInput{
				Bucket: &bucketName,
				Key:    &key,
				Body:   reader,
			}
			manager.NewUploader(awsClient)
			uploader := manager.NewUploader(awsClient)
			_, err := uploader.Upload(context.TODO(), &uploadInput)
			if err != nil {
				log.Fatalf("failed to upload key %s", key)
			}
			output <- fmt.Sprintf("Processed %s", key)
		}
		close(output)
	}()
	return output

}

// Issue is defined by integer Publisher ID, Publication ID and Issue ID.
type Issue struct {
	Publisher   int
	Publication int
	Issue       int
}

// FakeWriterAt implements a fake WriteAt method for synchronous writes, so we would be able to read/write S3 files from/to memory
type FakeWriterAt struct {
	w io.Writer
}

// WriteAt method ignores the *offset* because we forced sequential downloads
func (fw FakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	return fw.w.Write(p)
}
