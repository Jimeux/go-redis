package app

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func AWSConfig() aws.Config {
	cred := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("fake", "fake", ""))
	endpointResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		switch service {
		case dynamodb.ServiceID:
			return aws.Endpoint{
				PartitionID: "aws",
				URL:         "http://localhost:28000",
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{} // fallback to default
	})
	awsconf, _ := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(cred),
		config.WithEndpointResolver(endpointResolver),
		config.WithDefaultRegion("ap-northeast-1"),
	)
	return awsconf
}
