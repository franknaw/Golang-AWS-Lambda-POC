# testx
Proof of Concept for Go and AWS Lambda

#### Prerequisites
1. Install Golang
2. Instal AWS CLI
3. AWS Go SDK - go get github.com/aws/aws-lambda-go
4. AWS Lambda SDK - go get github.com/aws/aws-lambda-go/lambda
4. Excel SDK - go get github.com/360EntSecGroup-Skylar/excelize

#### Setup IAM role for lambda function
`{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "lambda.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}`

`aws iam create-role --role-name lambda-testx-executor \
--assume-role-policy-document file:///tmp/trust-policy.json`

Output as follows
`{
    "Role": {
        "RoleName": "lambda-testx-executor",
        "AssumeRolePolicyDocument": {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Principal": {
                        "Service": "lambda.amazonaws.com"
                    },
                    "Effect": "Allow",
                    "Action": "sts:AssumeRole"
                }
            ]
        },
        "RoleId": "SOME-ID",
        "Path": "/",
        "Arn": "arn:aws:iam::SOME-ID:role/lambda-testx-executor",
        "CreateDate": "2019-09-27T08:14:03Z"
    }
}`

Copy ARN (Amazon Resource Name) from results - for create-function listed below

#### Specify role permissions
`aws iam attach-role-policy --role-name lambda-testx-executor \
--policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole`

#### Add S3 role policy
`aws iam put-role-policy --role-name lambda-testx-executor \
--policy-name s3-item-crud-role \
--policy-document file:///tmp/privilege-policy.json`

`{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListAllMyBuckets",
                "s3:GetBucketLocation",
                "s3:GetObject",
                "s3:PutObject",
                "s3:ListBucket",
                "s3:DeleteBucket",
                "s3:CreateBucket"
            ],
            "Resource": "*"
        }
    ]
}`

#### Build Lambda Go program
`env GOOS=linux GOARCH=amd64 go build -o /tmp/main lambda && zip -j /tmp/main.zip /tmp/main`

#### Lambda Create Function
`aws lambda create-function --function-name testx --runtime go1.x \
--role arn:aws:iam::SOME-ID:role/lambda-testx-executor \
--handler main --zip-file fileb:///tmp/main.zip`

#### Invoke Lambda Service
`aws lambda invoke --function-name testx /tmp/output.json`

Output is as follows
```json
{
  "treasure-map.xlsx": "Excel File",
  "treasure1": "The Tollow, Anno 1401 in Germany, Störtebeker’s Golden Grave",
  "treasure10": "Isla Robinson Crusoe, Anno 1715 in Chile",
  "treasure2": "Gardiners Island, Anno 1699 in New York State, USA",
  "treasure3": "Oak Island, Anno 1795 in Nova Scotia, Canada, The Unknown Treasure",
  "treasure4": "Isla de Coco, 1820 in Panama, The famous treasure of Lima on Cocos Island",
  "treasure5": "Norman Island, Anno 1883 in the British Virgin Islands, Stevenson’s Treasure Island - The Ultimate Legend",
  "treasure6": "Ailsa Craig, Anno in the 17th century in Scotland",
  "treasure7": "Frégate Island, Anno 1730 in the Seychelles",
  "treasure8": "St. Joseph Atoll, Anno 1721 in the Seychelles, Concerning the Treasure of the St. Joseph Atoll",
  "treasure9": "Tupai, Anno 1822 in French Polynesia"
}
```
#### Re-deploy Lambda Service
`aws lambda update-function-code --function-name testx \
--zip-file fileb:///tmp/main.zip`

#### Delete Lambda Service - when needed
`aws lambda delete-function --function-name testx`

#### Query Cloudwatch logs
`aws logs filter-log-events --log-group-name /aws/lambda/testx \
--filter-pattern "ERROR"`

#### TODO
- Create Serverless Gateway/REST API
- Create unit tests and document code
- Create Terraform code to deploy Lambda function
- Create CircleCI build pipeline
- Cron Lambda function


