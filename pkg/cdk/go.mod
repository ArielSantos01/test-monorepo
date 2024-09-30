module stackcdk

go 1.22

replace (
	apigateway => ./../../infra/cdk/apigateway
	build => ./../../infra/cdk/build
	cdkcron => ./../../infra/cdk/cron
	cdklhttp => ./../../infra/cdk/lhttp
	cdksqs => ./../../infra/cdk/sqs
)

require (
	github.com/aws/aws-cdk-go/awscdk/v2 v2.159.1
	github.com/aws/constructs-go/constructs/v10 v10.3.0
)

require (
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/aws/jsii-runtime-go v1.103.1 // indirect
	github.com/cdklabs/awscdk-asset-awscli-go/awscliv1/v2 v2.2.202 // indirect
	github.com/cdklabs/awscdk-asset-kubectl-go/kubectlv20/v2 v2.1.2 // indirect
	github.com/cdklabs/awscdk-asset-node-proxy-agent-go/nodeproxyagentv6/v2 v2.1.0 // indirect
	github.com/cdklabs/cloud-assembly-schema-go/awscdkcloudassemblyschema/v36 v36.0.24 // indirect
)
