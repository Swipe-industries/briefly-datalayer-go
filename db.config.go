package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DBClient will be exported and used in other files
var DBClient *dynamodb.Client

func InitDB() {
	// Load AWS default config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Initialize DynamoDB client
	DBClient = dynamodb.NewFromConfig(cfg)
	fmt.Println("✅ DynamoDB connection established successfully")
}
