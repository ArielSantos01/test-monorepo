package apigateway

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"pkl"
	"strconv"
)

type Route struct {
	Method string `json:"method"`
	Route  string `json:"route"`
}
type Routes struct {
	Routes []Route `json:"routes"`
	ApiId  string  `json:"apiId"`
}

func CreateFunction(stack constructs.Construct, name, path string, opts ...Option) error {
	functionOpts := options{}
	for _, opt := range opts {
		opt(&functionOpts)
	}

	apiFunction := awslambda.NewFunction(stack, jsii.String(name+"function"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("./.bin/"+path+"/bootstrap.zip"), nil),
		Architecture: awslambda.Architecture_ARM_64(),
		MemorySize:   jsii.Number(2048),
	})

	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("LambdaIntegration"),
		apiFunction,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	routes, err := pkl.ReadConfig[Routes](path, functionOpts.pklConfig)
	if err != nil {
		return err
	}

	httpApiGateway := awsapigatewayv2.HttpApi_FromHttpApiAttributes(
		stack,
		jsii.String(name+routes.ApiId),
		&awsapigatewayv2.HttpApiAttributes{
			HttpApiId: jsii.String(routes.ApiId),
		})

	if err = createRoute(stack, httpApiGateway, integration, routes, name); err != nil {
		return err
	}
	return nil
}

func createRoute(stack constructs.Construct, httpApiGateway awsapigatewayv2.IHttpApi,
	integration awsapigatewayv2integrations.HttpLambdaIntegration, routes Routes, name string) error {

	for i, route := range routes.Routes {
		awsapigatewayv2.NewHttpRoute(stack, jsii.String(name+"-Route-"+strconv.Itoa(i)), &awsapigatewayv2.HttpRouteProps{
			HttpApi:     httpApiGateway,
			RouteKey:    awsapigatewayv2.HttpRouteKey_With(jsii.String(route.Route), awsapigatewayv2.HttpMethod(route.Method)),
			Integration: integration,
		})
	}

	return nil
}
