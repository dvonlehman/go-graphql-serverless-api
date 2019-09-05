SAM_TEMPLATE=sam.yaml
SAM_PACKAGED=sam-packaged.yaml
LAMBDA_NAME=go-graphql-api
MODULE=example.com/go-graphql-api
BINARY_NAME=go-graphql-api
S3_BUCKET=david-lambda-deployments
ZIP_FILE=dist.zip
    
all: clean test build package

build: clean
	GOOS=linux go build -o ./bin/$(BINARY_NAME)

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

update: build package
	aws lambda update-function-code --function-name $(LAMBDA_NAME) --zip-file fileb://$(ZIP_FILE) --publish
