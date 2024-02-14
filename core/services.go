package core

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// map of regions

type services struct {
	config      aws.Config
	autoscaling *autoscaling.Client
	ec2         *ec2.Client
}

type globalServices struct {
	config         aws.Config
	cloudformation *cloudformation.Client
	s3             *s3.Client
}

// List Stacks example
// res, err := cfn.ListStacks(context.TODO(),
// 	&cloudformation.ListStacksInput{})
// if err != nil {
// 	log.Println("Error while listing CloudFormation stacks: ", err.Error())
// 	return
// }
// for _, r := range res.StackSummaries {
// 	log.Println(*r.StackName, r.CreationTime)
// }
