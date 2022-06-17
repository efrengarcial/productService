package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/apex/gateway/v2"
	"github.com/apps/productService/products/delivery"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
	log "github.com/sirupsen/logrus"
	"gocloud.dev/docstore/awsdynamodb"

	"github.com/apps/productService/products"
)

var logger = log.New()

func main() {
	port := os.Getenv("PORT")
	mode := os.Getenv("MODE")
	logger.Out = os.Stdout
	logger.Level = log.InfoLevel
	if mode == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		log.SetFormatter(&log.TextFormatter{})
	}

	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	logger.Info("=======================================")
	logger.Info("Runinng gin-lambda server in " + addr)
	logger.Info("=======================================")


	usecase, cleanup, err := setup(mode != "production")
	if err != nil {
		logger.Fatal(err)
	}
	defer cleanup()

	router, err := delivery.Routes(usecase, logger)
	if err != nil {
		logger.Panic(err)
	}

	if mode == "production" {
		logger.Fatal(gateway.ListenAndServe(addr, router))
	} else {
		logger.Fatal(http.ListenAndServe(addr, router))
	}
}


// Init sets up an instance of this domains
// usecase, pre-configured with the dependencies.
func setup(integration bool) (usecase products.ProductService,cleanup func(), err error) {

	addCleanup := func(f func()) {
		old := cleanup
		cleanup = func() { old(); f() }
	}

	defer func() {
		if err != nil {
			cleanup()
			cleanup = nil
		}
	}()

	cleanup = func() {}

	region := os.Getenv("AWS_REGION")
	endPoint := os.Getenv("AWS_DYNAMODB_ENDPOINT")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region), Endpoint: aws.String(endPoint)},
	)
	if err != nil {
		return nil, cleanup, err
	}

	tableName := os.Getenv("TABLE_NAME")
	ddb := dynamodb.New(sess)
	coll, err := awsdynamodb.OpenCollection(
		ddb, tableName, "id", "", nil)
	if err != nil {
		return nil, cleanup, err
	}
	addCleanup(func() {
		fmt.Println("coll.Close()", coll.Close())
	})

	if integration == false {
		xray.Configure(xray.Config{LogLevel: "trace"})
		xray.AWS(ddb.Client)
	}

	//tableName := os.Getenv("TABLE_NAME")
	repository := products.NewDynamoDBRepository(coll)
	usecase = &products.Usecase{
		Repository: repository,
	}
	return usecase, cleanup, nil
}
