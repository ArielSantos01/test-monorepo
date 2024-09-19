package cdkcron

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func CreateFunction(stack constructs.Construct, name string, path string) awslambda.Function {

	cronFunction := awslambda.NewFunction(stack, jsii.String(name+"-cron"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("./.bin/"+path+"/bootstrap.zip"), nil),
		Architecture: awslambda.Architecture_ARM_64(),
	})

	rule := awsevents.NewRule(stack, jsii.String(name+"Rule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Expression(jsii.String("rate(2 minutes)")),
	})

	rule.AddTarget(awseventstargets.NewLambdaFunction(cronFunction, &awseventstargets.LambdaFunctionProps{}))

	// Outputs
	awscdk.NewCfnOutput(stack, jsii.String(name+"-CronLambdaArn"), &awscdk.CfnOutputProps{
		Value:       cronFunction.FunctionArn(),
		Description: jsii.String("The ARN of the Cron Lambda function"),
	})

	return cronFunction
}
