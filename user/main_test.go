package main

import (
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextRandom() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

func TestGetUser(t *testing.T) {
	ep := os.Getenv("DYNAMODB_ENDPOINT")
	r := os.Getenv("AWS_REGION")
	if len(r) == 0 {
		r = "ap-northeast-1"
	}

	tn := os.Getenv("DYNAMODB_TABLE_NAME")
	tn = tn + nextRandom()

	dyn := dynamodb.New(
		session.Must(
			session.NewSession(
				&aws.Config{
					Endpoint: aws.String(ep),
					Region:   aws.String(r),
				},
			),
		),
	)

	if err := deleteTestData(dyn, tn); err != nil {
		t.Logf("delete table failed %v\n", err)
	}

	if err := createTestData(dyn, tn); err != nil {
		t.Fatalf("table created failed rr: %v\n", err)
	}

	h := generateHandler(ep, r, tn)
	in := events.APIGatewayProxyRequest{}
	in.PathParameters = make(map[string]string)
	in.PathParameters["id"] = "1"

	res, err := h(in)
	if err != nil {
		t.Fatalf("err: %v\n", err)
	}
	t.Logf("result = %+v\n", res)
}

func createTestData(db *dynamodb.DynamoDB, tn string) error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("N"),
			},
		},
		TableName: aws.String(tn),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(2),
			WriteCapacityUnits: aws.Int64(2),
		},
	}

	if _, err := db.CreateTable(input); err != nil {
		return err
	}

	// PutItem
	putParams := &dynamodb.PutItemInput{
		TableName: aws.String(tn),
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String("1"),
			},
			"Name": {
				S: aws.String("Alice"),
			},
		},
	}

	if _, err := db.PutItem(putParams); err != nil {
		return err
	}
	return nil
}
func deleteTestData(db *dynamodb.DynamoDB, tn string) error {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(tn),
	}
	_, err := db.DeleteTable(input)
	if err != nil {
		return err
	}
	return nil
}
