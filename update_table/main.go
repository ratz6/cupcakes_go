package main

import (
	"math/rand"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func randgen(start int, end int, count int) []int {

	nums := make([]int, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for len(nums) < count {
		num := r.Intn(end-start) + start
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}

func Handler() {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Update item in table
	tableName := "cupcakes-go"
	res := randgen(1, 207, 100)
	currentTime := time.Now()
	user_name := make([]string, 0)
	user_name = append(user_name,
	"Alice",
	"Bob",
	"Charlie",
	"Denver" )
	rand.Seed(time.Now().Unix())
	for i := 0; i < len(res); i++ {
		Id := res[i]
		Cupcake := rand.Intn(101)
		Update_time := currentTime.Format("2006.01.02 15:04:05")
		Updated_by := user_name[rand.Intn(len(user_name))]
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":u": {
					S: aws.String(Update_time),
				},
				":c": {
					N: aws.String(strconv.Itoa(Cupcake)),
				},
				":n": {
					S: aws.String(Updated_by),
				},

			},
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"Id": {
					N: aws.String(strconv.Itoa(Id)),
				},
			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set Update_time = :u, Cupcake = :c, Updated_by = :n"),
		}

		_, err := svc.UpdateItem(input)
		if err != nil {
			log.Fatalf("Got error calling UpdateItem: %s", err)
		}
	}
	fmt.Println("Successfully updated ")
}

func main() {
	lambda.Start(Handler)
}
