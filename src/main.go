package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ignoring missing .env file")
	}
}

func main() {
	awsConfig := getAWSConfig()
	ociConfig := getOCIConfig()
	prefix := "9/17/21650"
	sync(awsConfig, ociConfig, prefix)
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

func sync(awsConfig aws.Config, ociConfig aws.Config, prefix string) {
	buckets := strings.Split(os.Getenv("BUCKETS"), ",")
	for _, bucketName := range buckets {
		iterObjects(bucketName, ociConfig, prefix)
		iterObjects(bucketName, awsConfig, prefix)
	}
}

func iterObjects(bucketName string, config aws.Config, prefix string) {
	// filter := "editoras/278/titulos/548/edicoes"
	fmt.Println("Object:", prefix)
	client := s3.NewFromConfig(config)
	params := &s3.ListObjectsV2Input{Bucket: &bucketName, Prefix: &prefix}
	paginator := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		o.Limit = 5
	})

	var i int
	log.Println("Objects: ", bucketName)
	for paginator.HasMorePages() {
		i++

		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		// Log the objects found
		for _, obj := range page.Contents {
			fmt.Println("Object:", *obj.Key)
		}

		if i > 0 {
			break
		}
	}
}
