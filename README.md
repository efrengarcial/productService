AWS_REGION=us-east-1;TABLE_NAME=products-integration;AWS_DYNAMODB_ENDPOINT=http://localhost:8000

aws cloudformation create-stack --stack-name  ecommerce-stack --template-body file://infrastructure/datastore.yml

aws cloudformation   update-stack --stack-name  ecommerce-stack --template-body file://infrastructure/datastore.yml
