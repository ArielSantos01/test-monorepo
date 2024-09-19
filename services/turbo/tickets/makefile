#bootstrap:
#	cdk bootstrap aws://776658659836/us-east-2 --profile draftea-dev -v
LAMBDAS := $(shell find cmd -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)

install:
	echo "Installing CDK"
	npm install -g aws-cdk

synt:
	echo "Synth CDK stack"
	cd cdk && cdk synth

build:
	LAMBDAS := $(shell find cmd -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)

	echo "Building CDK stack"
	@for lambda in $(LAMBDAS); do \
		rm -rf cdk/.bin/$$lambda && \
		GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -tags lambda.norpc -o cdk/.bin/$$lambda/bootstrap cmd/$$lambda/main.go && \
		zip -j cdk/.bin/$$lambda/bootstrap.zip cdk/.bin/$$lambda/bootstrap; \
	done

deploy: build
	echo "Deploying CDK stack"
	cd cdk && cdk deploy -v -y --require-approval never
	cd -

destroy:
	cd cdk && cdk destroy -v
	cd -