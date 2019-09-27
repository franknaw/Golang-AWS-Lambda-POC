package main

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
)

const (
	S3Region  = "us-east-1"
	S3Bucket  = "fac-cor"
	ExcelFile = "treasure-map.xlsx"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func run() (string, error) {

	s, err := session.NewSession(&aws.Config{Region: aws.String(S3Region)})
	if err != nil {
		return "Unable to create session", err
	}
	treasureMap := make(map[string]string)
	_ = getS3Objects(s, treasureMap)
	putObject(s, buildExcel(treasureMap))

	resultMap := make(map[string]string)
	listBucketObjects(s, resultMap)
	jsonOutput, _ := json.Marshal(resultMap)
	return string(jsonOutput), err
}

func buildExcel(treasureMap map[string]string) *excelize.File {
	var excelFile = excelize.NewFile()
	for k, v := range treasureMap {
		excelFile.NewSheet(k)
		_ = excelFile.SetColWidth(k, "A", "A", 350)
		_ = excelFile.SetCellValue(k, "A1", v)
	}

	return excelFile
}

func serviceError(err error) {
	errorLogger.Println(err.Error())
}

func main() {
	lambda.Start(run)
}
