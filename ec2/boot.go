package main

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func startEC2(instanceID string, dryRun bool) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess)

	instantStop := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		DryRun: aws.Bool(dryRun),
	}

	_, err := svc.StartInstances(instantStop)
	if err != nil {
		return err
	}

	return nil
}

func run() error {
	instanceID := os.Getenv("INSTANCE_ID")
	dryRun := strings.ToLower(os.Getenv("DRY_RUN")) == "true"
	return startEC2(instanceID, dryRun)
}

func main() {
	lambda.Start(run)
}
