package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	S3_REGION     = "us-east-1"
	S3_BUCKET     = "fac-cor"
	S3_BUCKET_KEY = "some-test.txt"
	S3_PAYLOAD    = "Some Test"
)

var treasureIslands = map[string]string{
	"treasure1":  "The Tollow, Anno 1401 in Germany, Störtebeker’s Golden Grave",
	"treasure2":  "Gardiners Island, Anno 1699 in New York State, USA",
	"treasure3":  "Oak Island, Anno 1795 in Nova Scotia, Canada, The Unknown Treasure",
	"treasure4":  "Isla de Coco, 1820 in Panama, The famous treasure of Lima on Cocos Island",
	"treasure5":  "Norman Island, Anno 1883 in the British Virgin Islands, Stevenson’s Treasure Island - The Ultimate Legend",
	"treasure6":  "Ailsa Craig, Anno in the 17th century in Scotland",
	"treasure7":  "Frégate Island, Anno 1730 in the Seychelles",
	"treasure8":  "St. Joseph Atoll, Anno 1721 in the Seychelles, Concerning the Treasure of the St. Joseph Atoll",
	"treasure9":  "Tupai, Anno 1822 in French Polynesia",
	"treasure10": "Isla Robinson Crusoe, Anno 1715 in Chile",
}

func s3Operation(item string) (string, error) {

	s, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})
	if err != nil {
		exitErrorf("Unable to establish session, %s", err)
	}

	buf := bytes.Buffer{}

	switch item {
	case "listBuckets":
		listBuckets(s)
	case "addItem":
		addObject(s, S3_BUCKET_KEY, S3_PAYLOAD)
	case "listObjects":
		listObjects(s)
	case "getObject":
		getBucketObject(s)
	case "createBucket":
		createBucket(s)
	case "deleteBucket":
		deleteBucket(s)
	case "deleteObject":
		deleteObject(s)
	case "loadTest":
		loadTest(s)
	case "deleteTest":
		deleteTest(s)
	case "deleteTestFile":
		deleteTestFile(s, "treasure-map.xlsx")
	}

	return buf.String(), err
}

func loadTest(s *session.Session) {
	for k, v := range treasureIslands {
		fmt.Printf("Adding %s -> %s\n", k, v)
		addObject(s, k, v)
	}
}

func deleteTest(s *session.Session) {
	for k, v := range treasureIslands {
		req, resp := s3.New(s).DeleteObjectRequest(&s3.DeleteObjectInput{
			Bucket: aws.String(S3_BUCKET),
			Key:    aws.String(k)})
		_ = req.Send()
		fmt.Printf("Deleting %s -> %s\n", resp.String(), v)
	}
}

func deleteTestFile(s *session.Session, key string) {
		req, _ := s3.New(s).DeleteObjectRequest(&s3.DeleteObjectInput{
			Bucket: aws.String(S3_BUCKET),
			Key:    aws.String(key)})
		_ = req.Send()
		fmt.Printf("Deleting %s\n", key)
}


func listObjects(s *session.Session) {

	err := s3.New(s).ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(S3_BUCKET),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			fmt.Println("S3 Object: ", *obj.Key)
		}
		return true
	})

	if err != nil {
		exitErrorf("Unable to list S3 Object, %s", err)
	}
}

func listBuckets(s *session.Session) {

	result, err := s3.New(s).ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list S3 Buckets, %s", err)
	}

	for _, b := range result.Buckets {
		fmt.Println("Bucket: ", aws.StringValue(b.Name)+"  ..created on: ", aws.TimeValue(b.CreationDate))
	}
}

func getBucketObject(s *session.Session) {
	resp, err := s3.New(s).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(S3_BUCKET),
		Key:    aws.String(S3_BUCKET_KEY),
	})

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		exitErrorf("Unable to get object, %s", err)
	}

	fmt.Println("S3 Bucket Object: ", string(body))
}

func createBucket(s *session.Session) {
	_, err := s3.New(s).CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(S3_BUCKET),
	})

	if err != nil {
		exitErrorf("Unable to create bucket, %s", err)
	}

	setBucketPolicy(s)
	setBucketACL(s)
}
func setBucketPolicy(s *session.Session) {

	readWriteAnonUserPolicy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Sid":       "Allow",
				"Effect":    "Allow",
				"Principal": "*",
				"Action": []string{
					"s3:GetObject",
				},
				"Resource": []string{
					fmt.Sprintf("arn:aws:s3:::%s/*", S3_BUCKET),
				},
			},
		},
	}

	policy, err := json.Marshal(readWriteAnonUserPolicy)

	_, err = s3.New(s).PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(S3_BUCKET),
		Policy: aws.String(string(policy)),
	})

	if err != nil {
		exitErrorf("Unable to set bucket policy, %s", err)
	}

	fmt.Printf("Successfully set bucket %q's policy\n", S3_BUCKET)

}

func setBucketACL(s *session.Session) {

	input := &s3.PutBucketAclInput{
		Bucket:           aws.String(S3_BUCKET),
		GrantFullControl: aws.String("id=someid"),
	}

	result, err := s3.New(s).PutBucketAcl(input)

	if err != nil {
		exitErrorf("Unable to set bucket ACL, %s", err)
	}

	fmt.Println("Successfully set bucket ACL", result)

}

func deleteBucket(s *session.Session) {
	_, err := s3.New(s).DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(S3_BUCKET),
	})

	if err != nil {
		exitErrorf("Unable to delete bucket, %s", err)
	}
}

func deleteObject(s *session.Session) {
	req, resp := s3.New(s).DeleteObjectRequest(&s3.DeleteObjectInput{
		Bucket: aws.String(S3_BUCKET),
		Key:    aws.String(S3_BUCKET_KEY)})

	err := req.Send()

	if err != nil {
		exitErrorf("Unable to delete object, %s", err)
	}

	fmt.Println("S3 object deleted: ", resp.String())
}

func addObject(s *session.Session, key string, payload string) {

	var err error
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(S3_BUCKET),
		Key:                  aws.String(key),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader([]byte(payload)),
		ContentLength:        aws.Int64(int64(len([]byte(payload)))),
		ContentType:          aws.String(http.DetectContentType([]byte(payload))),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		exitErrorf("Unable to add object, %s", err)
	}
}

func exitErrorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
	os.Exit(1)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("ff")
		exitErrorf("Please enter S3 Operation, %s", os.Args[0])
	}

	fmt.Println(os.Args[1])

	results, err := s3Operation(os.Args[1])

	if err != nil {
		exitErrorf("Error for item %q, %v", results, err)
	}

	fmt.Println(results)

}
