.PHONY: build run dev test lint clean deps

SAM_TEMPLATE=sam.yaml
SAM_PACKAGED=sam-packaged.yaml
LAMBDA_NAME=go-graphql-api
MODULE=example.com/go-graphql-api
BINARY_NAME=go-graphql-api
S3_BUCKET=david-lambda-deployments
ZIP_FILE=dist.zip
PORT=6000
    
all: clean test build package

build: clean
	go build -o ./bin/$(BINARY_NAME)

deps:
	go mod tidy

package: build
	(cd bin; zip ../$(ZIP_FILE) $(BINARY_NAME))

test:
	go test $(MODULE)/api -v

clean:
	rm -rf $(ZIP_FILE)
	rm -rf ./bin

provision: package
	sam package --template-file $(SAM_TEMPLATE) --output-template-file $(SAM_PACKAGED) --s3-bucket $(S3_BUCKET)
	sam deploy --template-file $(SAM_PACKAGED) --stack-name $(LAMBDA_NAME) --capabilities CAPABILITY_IAM

update: export GOOS = linux
update: build package
	aws lambda update-function-code --function-name $(LAMBDA_NAME) --zip-file fileb://$(ZIP_FILE) --publish

run: build
	./bin/$(BINARY_NAME) local

# Using reflex to watch for changes to .go file
# and re-run `make run`
# https://github.com/cespare/reflex/issues/50#issuecomment-388099690
# Install with `go get github.com/cespare/reflex`
local:
	LOG_LEVEL=debug reflex --start-service -r '\.go$$' make run
