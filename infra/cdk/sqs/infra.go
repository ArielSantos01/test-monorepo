package cdksqs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func CreateFunction(stack constructs.Construct, name string, path string) awslambda.Function {

	sqsFunction := awslambda.NewFunction(stack, jsii.String(name+"function"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("./.bin/"+path+"/bootstrap.zip"), nil),
		Architecture: awslambda.Architecture_ARM_64(),
	})
	// todo: add queueArn as parameter

	queueArn := "arn:aws:sqs:us-east-1:776658659836:QueeCdk"
	nameQueue := "QueueCdk"
	queue := awssqs.Queue_FromQueueArn(stack, jsii.String(nameQueue), jsii.String(queueArn))
	sqsFunction.AddEventSource(awslambdaeventsources.NewSqsEventSource(queue, &awslambdaeventsources.SqsEventSourceProps{}))
	queue.GrantConsumeMessages(sqsFunction)

	return sqsFunction
}
