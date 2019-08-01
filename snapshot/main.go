package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func stopEC2(instanceID string, dryRun bool) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess)

	instantStop := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		DryRun: aws.Bool(dryRun),
	}

	_, err := svc.StopInstances(instantStop)
	if err != nil {
		return err
	}

	return nil
}

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

func takeSnapShot(volumeID string, dryRun bool) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := ec2.New(sess)

	createSnapshot := &ec2.CreateSnapshotInput{
		VolumeId: aws.String(volumeID),
		DryRun:   aws.Bool(dryRun),
	}
	req, res := svc.CreateSnapshotRequest(createSnapshot)
	err := req.Send()
	if err != nil {
		return err
	}
	log.Println(res.String())

	return nil
}

func shutdownTakeSnapshot() (string, error) {
	instanceID := os.Getenv("INSTANCE_ID")
	volumeID := os.Getenv("VOLUME_ID")
	dryRun := strings.ToLower(os.Getenv("DRY_RUN")) == "true"

	err := stopEC2(instanceID, dryRun)
	if err != nil {
		log.Println(err)
	}

	if err = takeSnapShot(volumeID, dryRun); err != nil {
		log.Println(err)
	}

	// continue start ec2.
	time.Sleep(1 * time.Minute)
	err = startEC2(instanceID, dryRun)

	return "SUCCESS", err
}

func main() {
	lambda.Start(shutdownTakeSnapshot)
}
