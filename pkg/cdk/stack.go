package stackcdk

import (
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

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) (awscdk.Stack, error) {
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
	if sqsMap, exists := pathsMaps["sqs"]; exists {
		createSQSLambdas(stack, sqsMap)
	}

	if cronMap, exists := pathsMaps["cron"]; exists {
		createCRONLambdas(stack, cronMap)
	}

	return stack, nil
}

func createLHTTPLambdas(stack constructs.Construct, httpMap map[string]string) {
	for nameFunc, path := range httpMap {
		cdklhttp.CreateFunction(stack, nameFunc, path)
	}
}

func createSQSLambdas(stack constructs.Construct, sqsMap map[string]string) {
	for nameFunc, path := range sqsMap {
		cdksqs.CreateFunction(stack, nameFunc, path)
	}
}

func createCRONLambdas(stack constructs.Construct, cronMap map[string]string) {
	for nameFunc, path := range cronMap {
		cdkcron.CreateFunction(stack, nameFunc, path)
	}
}

func BuildFunc() error {

	if err := build.Exec(); err != nil {
		return err
	}
	return nil
}
