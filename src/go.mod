module s3sync

go 1.16

require (
	copier v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go-v2 v1.6.0
	github.com/aws/aws-sdk-go-v2/config v1.3.0
	github.com/aws/aws-sdk-go-v2/credentials v1.2.1
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.2.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.10.0
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.3.0
)

replace copier => ./copier
