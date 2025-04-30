scaffold:
	@echo "Scaffolding..."
	@mkdir  -p auth-service/{cmd,internal}
	@mkdir -p auth-service/internal/{delivery,model, repo,config}
	@mkdir  -p auth-service/cmd/{server,grpc}
	@mkdir 	-p common/{proto,jwt-pkg}
build:
	@GOOS=linux GOARCH=arm64 go build -o email-function main.go

zip:
	@zip function.zip email-function

lambda-load:
	@aws lambda update-function-code \
	--function-name email-function \
	--zip-file fileb://function.zip \
	--region us-east-1

lambda-create:
	@aws lambda create-function \
       --function-name email-function \
       --runtime provided.al2 \
       --handler main \
       --role arn:aws:iam::162047532564:role/lambda-execution-role \
       --zip-file fileb://function.zip \
       --region us-east-1

ecr:
	@aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <aws_account_id>.dkr.ecr.<region>.amazonaws.com

login-ecr:
	@aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 162047532564.dkr.ecr.us-east-1.amazonaws.com

docker-node-build:
	@docker buildx build --platform=linux/amd64 -t swan_auth .

docker-node-tag:
	@docker tag swan-node:latest 162047532564.dkr.ecr.us-east-1.amazonaws.com/swan-node:latest

docker-node-push:
	@docker push 162047532564.dkr.ecr.us-east-1.amazonaws.com/swan-node:latest

ec-2-login:
	@ssh -i /Users/swanhtet1aungphyo/Downloads/swanhtet.pem ubuntu@ec2-54-167-62-240.compute-1.amazonaws.com

aws-account-id:
	@aws sts get-callerr-identtity

list:
	@aws lambda list-functions --region us-east-1


proto:
	@protoc --go_out=. --go-grpc_out=. \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       ./common/proto/user.proto