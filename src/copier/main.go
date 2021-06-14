package copier

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
)

var pageLimit int32 = 1000

// CopySource defines the bucket name and how to parse this prefix, given an issue
type CopySource struct {
	BucketName string
	PrefixFunc func(Issue) string
}

var buckets = []CopySource{
	CopySource{
		BucketName: "revistas.magtab.com",
		PrefixFunc: func(issue Issue) string {
			return fmt.Sprintf("%d/%d/%d", issue.Publisher, issue.Publication, issue.Issue)
		},
	},
	CopySource{
		BucketName: "revistas.magtab.com",
		PrefixFunc: func(issue Issue) string {
			return fmt.Sprintf("%d/%d/Titulo.json", issue.Publisher, issue.Publication)
		},
	},
}

// Copy files from OCI to AWS buckets
func Copy(issue Issue) {
	clients := Clients{
		AWS: getAWSClient(getAWSConfig()),
		OCI: getAWSClient(getOCIConfig()),
	}
	copyIssueFiles(clients, issue)
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

func getAWSClient(config aws.Config) *s3.Client {
	return s3.NewFromConfig(config)
}

func copyIssueFiles(clients Clients, issue Issue) {
	for _, item := range buckets {
		prefix := item.PrefixFunc(issue)
		ociItems := getObjectKeys(clients.OCI, item.BucketName, prefix)
		fmt.Println("=", item.BucketName)
		copier1 := copyObjects(clients, item.BucketName, ociItems)
		copier2 := copyObjects(clients, item.BucketName, ociItems)
		copier3 := copyObjects(clients, item.BucketName, ociItems)
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

func getObjectKeys(client *s3.Client, bucketName string, prefix string) <-chan string {
	params := &s3.ListObjectsV2Input{Bucket: &bucketName, Prefix: &prefix}
	paginator := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		o.Limit = pageLimit
	})
	channel := make(chan string)
	go func() {
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(context.TODO())
			if err != nil {
				log.Fatalf("failed to get page %v", err)
			}

			for _, obj := range page.Contents {
				channel <- *obj.Key
			}
		}
		close(channel)
	}()
	return channel
}

func copyObjects(clients Clients, bucketName string, input <-chan string) <-chan string {
	output := make(chan string)
	go func() {
		for key := range input {
			reader, writer := io.Pipe()
			downloadFile(clients.OCI, writer, bucketName, key)
			uploadFile(clients.AWS, reader, bucketName, key)
			output <- fmt.Sprintf("Processed %s", key)
		}
		close(output)
	}()
	return output
}

func downloadFile(client *s3.Client, writer *io.PipeWriter, bucketName string, key string) {
	requestInput := s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}
	downloader := manager.NewDownloader(client)
	downloader.Concurrency = 1
	go func() {
		defer writer.Close()
		downloader.Download(context.TODO(), FakeWriterAt{writer}, &requestInput)
	}()
}

func uploadFile(client *s3.Client, reader *io.PipeReader, bucketName string, key string) {
	uploadInput := s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   reader,
	}
	uploader := manager.NewUploader(client)
	_, err := uploader.Upload(context.TODO(), &uploadInput)
	if err != nil {
		log.Fatalf("failed to upload key %s", key)
	}
}

// Clients group AWS and OCI clients
type Clients struct {
	AWS *s3.Client
	OCI *s3.Client
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
