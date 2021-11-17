package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil" 
	"log"
	"net/http"
	"os"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := http.Get("https://search-cupcakes-go-4thiwnjwwx3ck64ssr7fd6v3ai.ap-south-1.es.amazonaws.com/lambda-index/_doc/_search?size=10000&sort=Id")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	} 
	result := map[string]json.RawMessage{}
	err = json.Unmarshal((responseData), &result)
	if err != nil {
		fmt.Println(err)
	}
	var tmp = result["hits"]
	s, err := json.Marshal(tmp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(s))
	
	apiResponse:= events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type": "application/json",
		},
		Body: string(s),
	} 
	return apiResponse,nil
}
func main() {
	lambda.Start(handler)
}
