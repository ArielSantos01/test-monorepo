package cdksqs

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"pkl"
)

type QueueArn struct {
	Arn string `json:"queueArn"`
}

func CreateFunction(stack constructs.Construct, name, path string, opts ...Option) error {

	functionOpts := options{}
	for _, opt := range opts {
		opt(&functionOpts)
	}

	sqsFunction := awslambda.NewFunction(stack, jsii.String(name+"function"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("./.bin/"+path+"/bootstrap.zip"), nil),
		Architecture: awslambda.Architecture_ARM_64(),
	})

	queueArn, err := pkl.ReadConfig[QueueArn](path, functionOpts.pklConfig)
	if err != nil {
		return err
	}
	fmt.Println(queueArn.Arn)
	queue := awssqs.Queue_FromQueueArn(stack, jsii.String(queueArn.Arn), jsii.String(queueArn.Arn))
	sqsFunction.AddEventSource(awslambdaeventsources.NewSqsEventSource(queue, &awslambdaeventsources.SqsEventSourceProps{}))
	queue.GrantConsumeMessages(sqsFunction)

	return nil
}
