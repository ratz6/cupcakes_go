package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"strconv"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)
type Item struct {
	Id  int      `json:"Id"`
	Month string 	`json:"Month"`
	Cupcake	int		`json:"Cupcake"`
	Updated_by string	`json:"Updated_by`
}
// Get table items from JSON file
func getItems() []Item {
	raw, err := ioutil.ReadFile("./user_data.json")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var items []Item
	json.Unmarshal(raw, &items)
	return items
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)
	if err != nil {
		log.Fatalf("Error creating session: %s", err)
	}
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	tableName := "cupcakes-go"
	currentTime := time.Now()
	items := getItems()
	for _, item := range items {

		_, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			log.Fatalf("Got error marshalling map: %s", err)
		}
		Id_u := item.Id
		Cupcake_u := item.Cupcake
		Updated_by_u := item.Updated_by
		Update_time := currentTime.Format("2006.01.02 15:04:05")
		
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":u": {
					S: aws.String(Update_time),
				},
				":c": {
					N: aws.String(strconv.Itoa(Cupcake_u)),
				},
				":n": {
					S: aws.String(Updated_by_u),
				},
	
			},
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"Id": {
					N: aws.String(strconv.Itoa(Id_u)),
				},
			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set Update_time = :u, Cupcake = :c, Updated_by = :n"),
		}
		_, err = svc.UpdateItem(input)

		if err != nil {
			log.Fatalf("Got error calling UpdateItem: %s", err)
		}
	}
	fmt.Println("Successfully updated table");
}
