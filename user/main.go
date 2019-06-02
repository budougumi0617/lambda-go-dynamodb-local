package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// User :
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func generateHandler(endpoint, region, tableName string) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// Request
		id, _ := request.PathParameters["id"]

		// DynamoDB
		sess := session.Must(session.NewSession())
		config := aws.NewConfig().WithRegion(region)
		if len(endpoint) > 0 {
			config = config.WithEndpoint(endpoint)
		}

		db := dynamodb.New(sess, config)
		response, err := db.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"Id": {
					N: aws.String(string(id)),
				},
			},
			AttributesToGet: []*string{
				aws.String("Id"),
				aws.String("Name"),
			},
			ConsistentRead:         aws.Bool(true),
			ReturnConsumedCapacity: aws.String("NONE"),
		})
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		user := User{}
		err = dynamodbattribute.Unmarshal(&dynamodb.AttributeValue{M: response.Item}, &user)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		// Json
		bytes, err := json.Marshal(user)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			Body:       string(bytes),
			StatusCode: http.StatusOK,
		}, nil
	}
}

func main() {
	// Environment variables
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	if len(region) == 0 {
		region = "ap-northeast-1"
	}
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	h := generateHandler(endpoint, region, tableName)
	lambda.Start(h)
}
