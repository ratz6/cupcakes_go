package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
)

// Create struct to hold info about new item
type Item struct {
	Id          int
	Month       string
	Cupcake     int
	Update_time string
}

func main() {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	tableName := "cupcakes-go"

	csvfile, err := os.Open("multiTimeline.csv")
	r := csv.NewReader(csvfile)
	if err != nil {
		log.Fatal(err)
	}
	header, _ := r.Read()
	fmt.Printf(" Headers: %v", header)

	for i := 1; ; i = i + 1 {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		val, _ := strconv.Atoi(record[1])
		item := Item{
			Id:          i,
			Month:       record[0],
			Cupcake:     val,
			Update_time: "-",
		}

		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			log.Fatalf("Got error marshalling new item: %s", err)
		}
		// Create item in table

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			log.Fatalf("Got error calling PutItem: %s", err)
		}
	}
	fmt.Println("The Data is loaded in DynamoDb")
}
