package storage

import (
	"authz/config"
	"authz/types"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StorageService struct {
	client *mongo.Client
}

//go:generate moq -out storageservice_mock_test.go . IStorageService
type IStorageService interface {
	GetUserHash(string) (string, error)
	GetUserId(string) (string, error)
}

var StorageSvc IStorageService

func InitStorage() error {
	log.Println("Connecting to DB")
	storageSvc := &StorageService{}
	uname := config.Configs.DBConfig.User
	pass := config.Configs.DBConfig.Password
	host := config.Configs.DBConfig.Host
	appname := config.Configs.DBConfig.AppName
	uri := "mongodb+srv://" + uname + ":" + pass + "@" + host + "/?retryWrites=true&w=majority&appName=" + appname
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("error connection to db: %v", err)
	}
	storageSvc.client = client
	StorageSvc = storageSvc
	return nil
}

func (s *StorageService) GetUserHash(email string) (string, error) {
	coll := s.client.Database(config.Configs.DBConfig.DBName).Collection(config.Configs.DBConfig.UsersCollection)
	var result types.MongoUserDoc
	err := coll.FindOne(context.Background(), bson.D{{"user_details.email", email}}).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error finding the user: %v", err)
	}
	fmt.Println(result)
	return result.UserDetails.PasswordHash, nil
}

func (s *StorageService) GetUserId(email string) (string, error) {
	coll := s.client.Database(config.Configs.DBConfig.DBName).Collection(config.Configs.DBConfig.UsersCollection)
	var result types.MongoUserDoc
	err := coll.FindOne(context.Background(), bson.D{{"user_details.email", email}}).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error finding the user: %v", err)
	}
	fmt.Println(result)
	return result.UserId, nil
}
