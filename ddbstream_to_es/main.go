package main

import (
  "fmt"
  "net/http"
  "strings"
  "time"
  "encoding/json"
  "strconv"
  "bytes"
  

  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/signer/v4"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/aws/aws-lambda-go/events"

)

type Item struct {
	Id          int
	Month       string
	Cupcake     int
	Update_time string
}

func handler(e events.DynamoDBEvent) error {

	// Basic information for the Amazon Elasticsearch Service domain
	domain := "https://search-cupcakes-go-4thiwnjwwx3ck64ssr7fd6v3ai.ap-south-1.es.amazonaws.com" // e.g. https://my-domain.region.es.amazonaws.com
	index := "lambda-index"
	region := "ap-south-1" // e.g. us-east-1
	service := "es"

	// Get credentials from environment variables and create the AWS Signature Version 4 signer
	credentials := credentials.NewEnvCredentials()
	signer := v4.NewSigner(credentials)

	var item map[string]events.DynamoDBAttributeValue

	for _,v := range e.Records {
		fmt.Println(v)
		switch v.EventName {
			case "INSERT":
				fallthrough
			case "MODIFY":
				tableName := strings.Split(v.EventSourceArn, "/")[1]
				fmt.Printf("tableName:%v", tableName)

				item = v.Change.NewImage

				id,_ := item["Id"].Integer()
				cup,_ := item["Cupcake"].Integer()

				endpoint := domain + "/" + index + "/" + "_doc" + "/" + strconv.Itoa(int(id))

				new_item := Item {
					Id: int(id),
					Month : item["Month"].String(),
					Cupcake: int(cup),
					Update_time : item["Update_time"].String(),
				}
				fmt.Printf("item:%+v", new_item)
				// Marshalling the Data
				b, err := json.Marshal(new_item)
				if err != nil {
					fmt.Printf("error in marshalling, err:%v", err)
					continue
				}
				final_b := strings.NewReader(string(b))

				fmt.Printf("req:%v", string(b))

				// An HTTP client for sending the request
				client := &http.Client{}

				// Form the HTTP request
				req, err := http.NewRequest(http.MethodPut, endpoint,bytes.NewBuffer(b))
				if err != nil {
					fmt.Print(err)
				}
				// You can probably infer Content-Type programmatically, but here, we just say that it's JSON
				req.Header.Add("Content-Type", "application/json")

				// Sign the request, send it, and print the response
				signer.Sign(req,final_b,service, region, time.Now())
				resp, err := client.Do(req)
				if err != nil {
					fmt.Print(err)
				}
				fmt.Print(resp.Status + "\n")
			default :
		}
	}
	fmt.Println("ES Added")
	return nil
}
func main() {
	lambda.Start(handler)
}