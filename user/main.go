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

// User is user struct.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func generateHandler(endpoint, region, tableName string) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// Request
		id, _ := request.PathParameters["id"]

		// Connect DynamoDB
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
		if err := dynamodbattribute.Unmarshal(&dynamodb.AttributeValue{M: response.Item}, &user); err != nil {
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
	ep := os.Getenv("DYNAMODB_ENDPOINT")
	r := os.Getenv("AWS_REGION")
	if len(r) == 0 {
		r = "ap-northeast-1"
	}
	tn := os.Getenv("DYNAMODB_TABLE_NAME")

	h := generateHandler(ep, r, tn)
	lambda.Start(h)
}
