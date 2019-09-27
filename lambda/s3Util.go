package main

import (
	"bytes"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"net/http"
)

func getS3Objects(s *session.Session, treasureMap map[string]string) error {
	err := s3.New(s).ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(S3Bucket),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			treasureMap[*obj.Key] = getObject(s, *obj.Key)
		}
		return true
	})

	return err
}

func getObject(s *session.Session, key string) string {
	resp, err := s3.New(s).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(key),
	})

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		serviceError(err)
	}

	return string(body)
}

func putObject(s *session.Session, excelFile *excelize.File) {
	var err error
	buffer, err := excelFile.WriteToBuffer()
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(S3Bucket),
		Key:                  aws.String(ExcelFile),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer.Bytes()),
		ContentLength:        aws.Int64(int64(len(buffer.Bytes()))),
		ContentType:          aws.String(http.DetectContentType(buffer.Bytes())),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		serviceError(err)
	}
}

func listBucketObjects(s *session.Session, resultMap map[string]string) {
	err := s3.New(s).ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(S3Bucket),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			if *obj.Key == ExcelFile {
				resultMap[*obj.Key] = "Excel File"
			} else {
				resultMap[*obj.Key] = getObject(s, *obj.Key)
			}
		}
		return true
	})

	if err != nil {
		serviceError(err)
	}
}
