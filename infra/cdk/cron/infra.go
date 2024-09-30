package cdkcron

import (
	"encoding/json"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"os"
	"path/filepath"
)

type Cron struct {
	Minute  string `json:"minute"`
	Hour    string `json:"hour"`
	Day     string `json:"day"`
	Month   string `json:"month"`
	WeekDay string `json:"weekDay"`
	Year    string `json:"year"`
}
type Schedule struct {
	Cron       Cron   `json:"cron"`
	Expression string `json:"expression"`
}

func CreateFunction(stack constructs.Construct, name string, path string) error {

	cronFunction := awslambda.NewFunction(stack, jsii.String(name+"-cron"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("./.bin/"+path+"/bootstrap.zip"), nil),
		Architecture: awslambda.Architecture_ARM_64(),
	})

	config, err := readConfig(path)
	if err != nil {
		return err
	}

	var rule awsevents.Rule

	if config.Expression != "" {
		rule = awsevents.NewRule(stack, jsii.String(name+"Rule"), &awsevents.RuleProps{
			Schedule: awsevents.Schedule_Expression(jsii.String(config.Expression)),
		})
		return nil
	} else {
		cronOptions := &awsevents.CronOptions{
			Minute: jsii.String(normalizeField(config.Cron.Minute)),
			Hour:   jsii.String(normalizeField(config.Cron.Hour)),
			Month:  jsii.String(normalizeField(config.Cron.Month)),
			Year:   jsii.String(normalizeField(config.Cron.Year)),
		}

		if config.Cron.Day != "" {
			cronOptions.Day = jsii.String(normalizeField(config.Cron.Day))
		} else {
			cronOptions.WeekDay = jsii.String(normalizeField(config.Cron.Day))
		}

		rule = awsevents.NewRule(stack, jsii.String(name+"Rule"), &awsevents.RuleProps{
			Schedule: awsevents.Schedule_Cron(cronOptions),
		})
	}

	rule.AddTarget(awseventstargets.NewLambdaFunction(cronFunction, &awseventstargets.LambdaFunctionProps{}))

	awscdk.NewCfnOutput(stack, jsii.String(name+"-CronLambdaArn"), &awscdk.CfnOutputProps{
		Value:       cronFunction.FunctionArn(),
		Description: jsii.String("The ARN of the Cron Lambda function"),
	})

	return nil
}

func normalizeField(field string) string {
	if field == "" {
		return "*"
	}
	return field
}

func readConfig(path string) (Schedule, error) {

	workingDir, err := os.Getwd()
	if err != nil {
		return Schedule{}, err
	}
	configPath := filepath.Join(workingDir, "..", "cmd", path, "/config.json")

	var configData []byte
	configData, err = os.ReadFile(configPath)
	if err != nil {
		return Schedule{}, err
	}

	var config Schedule

	if err = json.Unmarshal(configData, &config); err != nil {
		return Schedule{}, err
	}
	return config, nil
}
