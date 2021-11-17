package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
    "fmt"
    "log"
	"bytes"
	"encoding/json"
)
type Item struct {
    Id   int
    Month  string
    Cupcake int
    Updated_by string
}
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))  

    svc := dynamodb.New(sess)
    tableName := "cupcakes-go"
	updatedBy := "Nishit"
    
    filt := expression.Name("Updated_by").Equal(expression.Value(updatedBy))

    proj := expression.NamesList(expression.Name("Id"), expression.Name("Month"), expression.Name("Cupcake"),expression.Name("Updated_by"))

    expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
    if err != nil {
        log.Fatalf("Got error building expression: %s", err)
    }
    params := &dynamodb.ScanInput{
        ExpressionAttributeNames:  expr.Names(),
        ExpressionAttributeValues: expr.Values(),
        FilterExpression:          expr.Filter(),
        ProjectionExpression:      expr.Projection(),
        TableName:                 aws.String(tableName),
    }
    result, err := svc.Scan(params)
    if err != nil {
        log.Fatalf("Query API call failed: %s", err)
    }
    numItems := 0
	res := make([]Item,0)
	
    for _, i := range result.Items {
        item := Item{}
        err = dynamodbattribute.UnmarshalMap(i, &item)

        if err != nil {
            log.Fatalf("Got error unmarshalling: %s", err)
        }
        if item.Updated_by == updatedBy {
            numItems++
			fl := Item {
				Id: item.Id,
				Month: item.Month,
				Cupcake: item.Cupcake,
				Updated_by : item.Updated_by,
			}
			res =  append(res,fl)
            fmt.Println("Name: ", item.Updated_by)
            fmt.Println("Id:", item.Id)
            fmt.Println("Month:",item.Month)
        } 
	}
	responseBodyBytes := new(bytes.Buffer)
	json.NewEncoder(responseBodyBytes).Encode(res)
	apiResponse := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{ 			
			"Access-Control-Allow-Origin": "*",
			"Content-Type": "application/json",
		},
		Body: string(responseBodyBytes.Bytes()),
	}
	fmt.Println("Found ", numItems, " updates made having ", updatedBy)
	return apiResponse, nil
}
func main(){
	lambda.Start(Handler)
}