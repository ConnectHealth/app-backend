package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

const (
	PatientsTableName = "patients"
)

type Patient struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func GetPatients() ([]Patient, error) {
	svc, err := getSvc()
	if err != nil {
		return nil, err
	}

	// setup(svc)

	params := &dynamodb.ScanInput{
		TableName: aws.String(PatientsTableName), // Required
		AttributesToGet: []*string{
			aws.String("id"),
			aws.String("firstName"),
			aws.String("lastName"),
		},
		ConsistentRead: aws.Bool(true),
		Select:         aws.String(dynamodb.SelectSpecificAttributes),
	}
	resp, err := svc.Scan(params)
	if err != nil {
		fmt.Println("oh no")
		return nil, err
	}

	var patients []Patient
	for _, v := range resp.Items {
		patient := Patient{}
		err = dynamodbattribute.UnmarshalMap(v, &patient)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}

	return patients, nil
}

func setup(svc *dynamodb.DynamoDB) error {
	// delete table if exists
	deleteTable(svc, PatientsTableName)

	err := createPatientsTable(svc)
	if err != nil {
		fmt.Println("oh no")
		return err
	}

	// insert some data
	patients := []Patient{
		Patient{"1", "John", "Smith"},
		Patient{"2", "Jane", "Smith"},
		Patient{"3", "Jack", "Jones"},
	}

	for _, p := range patients {
		err = putPatient(svc, p)
		if err != nil {
			fmt.Println("oh no")
			return err
		}
	}
	return nil
}

func getSvc() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "could not create AWS session")
	}
	// ap-southeast-2
	return dynamodb.New(sess, &aws.Config{Endpoint: aws.String("http://localhost:8000"), Region: aws.String("eu-west-1")}), nil
}

func putPatient(svc *dynamodb.DynamoDB, patient Patient) error {
	item, err := dynamodbattribute.ConvertToMap(patient)
	if err != nil {
		return errors.Wrap(err, "could not convert patient")
	}
	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(PatientsTableName),
	}
	_, err = svc.PutItem(params)
	return err
}

func createPatientsTable(svc *dynamodb.DynamoDB) error {
	params := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{ // Required
			{ // Required
				AttributeName: aws.String("id"), // Required
				AttributeType: aws.String("S"),  // Required
			},
			{ // Required
				AttributeName: aws.String("lastName"), // Required
				AttributeType: aws.String("S"),        // Required
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{ // Required
			{ // Required
				AttributeName: aws.String("id"),   // Required
				KeyType:       aws.String("HASH"), // Required

			}, {
				AttributeName: aws.String("lastName"), // Required
				KeyType:       aws.String("RANGE"),    // Required
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ // Required
			ReadCapacityUnits:  aws.Int64(1), // Required
			WriteCapacityUnits: aws.Int64(1), // Required
		},
		TableName: aws.String(PatientsTableName), // Required
	}
	_, err := svc.CreateTable(params)
	return err
}

func deleteTable(svc *dynamodb.DynamoDB, tableName string) error {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(PatientsTableName), // Required
	}
	_, err := svc.DeleteTable(params)
	return err
}
