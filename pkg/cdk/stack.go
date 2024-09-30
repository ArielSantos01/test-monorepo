package stackcdk

import (
	"apigateway"
	"build"
	"cdkcron"
	"cdklhttp"
	"cdksqs"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"strings"
)

type CdkStackProps struct {
	awscdk.StackProps
}

var pathsMaps = make(map[string]map[string]string)

func processPath() error {
	if err := BuildFunc(); err != nil {
		return err
	}
	for _, path := range build.FoldersTree {
		parts := strings.Split(path, "/")
		funcType := parts[0]
		funcName := parts[len(parts)-1]

		if _, exists := pathsMaps[funcType]; !exists {
			pathsMaps[funcType] = make(map[string]string)
		}
		pathsMaps[funcType][funcName] = path
	}
	return nil
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps, opts ...Option) (awscdk.Stack, error) {
	functionOpts := options{}

	for _, opt := range opts {
		opt(&functionOpts)
	}

	if err := processPath(); err != nil {
		return nil, err
	}

	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	if httpMap, exists := pathsMaps["http"]; exists {
		createLHTTPLambdas(stack, httpMap)
	}

	if cronMap, exists := pathsMaps["cron"]; exists {
		err := createCRONLambdas(stack, cronMap)
		if err != nil {
			return nil, err
		}
	}

	if sqsMap, exists := pathsMaps["sqs"]; exists {
		err := createSQSLambdas(stack, sqsMap, functionOpts)
		if err != nil {
			return nil, err
		}
	}

	if apiHttpMap, exists := pathsMaps["apihttp"]; exists {
		err := createApiHttpLambdas(stack, apiHttpMap, functionOpts)
		if err != nil {
			return nil, err
		}
	}

	return stack, nil
}
func createApiHttpLambdas(stack constructs.Construct, apiHttpMap map[string]string, opt options) error {
	for nameFunc, path := range apiHttpMap {
		err := apigateway.CreateFunction(stack, nameFunc, path, apigateway.WithPklConfig(opt.pklConfig))
		if err != nil {
			return err
		}
	}
	return nil
}

func createSQSLambdas(stack constructs.Construct, sqsMap map[string]string, opt options) error {
	for nameFunc, path := range sqsMap {
		err := cdksqs.CreateFunction(stack, nameFunc, path, cdksqs.WithPklConfig(opt.pklConfig))
		if err != nil {
			return err
		}
	}
	return nil
}

func createLHTTPLambdas(stack constructs.Construct, httpMap map[string]string) {
	for nameFunc, path := range httpMap {
		cdklhttp.CreateFunction(stack, nameFunc, path)
	}
}

func createCRONLambdas(stack constructs.Construct, cronMap map[string]string) error {
	for nameFunc, path := range cronMap {
		err := cdkcron.CreateFunction(stack, nameFunc, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func BuildFunc() error {

	if err := build.Exec(); err != nil {
		return err
	}
	return nil
}
