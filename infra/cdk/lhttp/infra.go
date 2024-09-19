package cdklhttp

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func CreateFunction(stack constructs.Construct, name string, path string) awslambda.Function {

	myFunction := awslambda.NewFunction(stack, jsii.String(name), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("./.bin/"+path+"/bootstrap.zip"), nil),
		Architecture: awslambda.Architecture_ARM_64(),
	})

	myFunctionUrl := myFunction.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	awscdk.NewCfnOutput(stack, jsii.String(name+"-lambda-url"), &awscdk.CfnOutputProps{
		Value: myFunctionUrl.Url(),
	})

	return myFunction
}
