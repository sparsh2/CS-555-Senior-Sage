package storage

import (
	"storage-service/config"
	"storage-service/types"
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

//go:generate moq -out storageservice_mock.go . IStorageService
type IStorageService interface {
	GetUserHash(string) (string, error)
	GetUserId(string) (string, error)
	GetUserData(string, []types.ResourceType, *types.GetDataResponse) error
	InsertUserDoc(*types.MongoUserDoc) error
	GetAclDoc(string) (*types.MongoAclsDoc, error)
	InsertAclDoc(*types.MongoAclsDoc) error
	WriteUserData(string, map[types.ResourceType]string, []types.ResourceType) error
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

func (s *StorageService) WriteUserData(UserId string, data map[types.ResourceType]string, resources []types.ResourceType) error {
	coll := s.client.Database(config.Configs.DBConfig.DBName).Collection(config.Configs.DBConfig.UsersCollection)
	updateVal := bson.D{}
	for _, res := range resources {
		// TODO: check if the resource is valid and log error
		updateVal = append(updateVal, bson.E{Key: "data." + types.ResourceToMongoField[res], Value: data[res]})
	}
	_, err := coll.UpdateOne(context.Background(), bson.D{bson.E{Key: "user_id", Value: UserId}}, bson.D{bson.E{Key: "$set", Value: updateVal}})
	if err != nil {
		return fmt.Errorf("error updating user data: %v", err)
	}
	return nil
}

func (s *StorageService) GetUserData(UserId string, resources []types.ResourceType, resp *types.GetDataResponse) error {
	coll := s.client.Database(config.Configs.DBConfig.DBName).Collection(config.Configs.DBConfig.UsersCollection)
	projection := bson.D{}
	for _, r := range resources {
		projection = append(projection, bson.E{Key: "data."+types.ResourceToMongoField[r], Value: 1})
	}
	var result types.MongoUserDoc
	err := coll.FindOne(context.Background(), bson.D{bson.E{Key: "user_id", Value: UserId}}, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return fmt.Errorf("error finding the user: %v", err)
	}
	for _, r := range resources {
		switch r {
		case types.RESOURCE_USER_DETAILS:
			resp.Data[r] = result.Data.UserDetails
		case types.RESOURCE_USER_PREFERENCES:
			resp.Data[r] = result.Data.UserPreferences
		case types.RESOURCE_USER_REMINDERS:
			resp.Data[r] = result.Data.UserReminders
		}
	}
	return nil
}

func (s *StorageService) GetAclDoc(UserId string) (*types.MongoAclsDoc, error) {
	coll := s.client.Database(config.Configs.DBConfig.DBName).Collection(config.Configs.DBConfig.AclsCollection)
	var result types.MongoAclsDoc
	err := coll.FindOne(context.Background(), bson.D{bson.E{Key: "uid", Value: UserId}}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error finding Acls for the user: %v", err)
	}
	return &result, nil
}

func (s *StorageService) GetUserHash(email string) (string, error) {
	coll := s.client.Database(config.Configs.DBConfig.DBName).Collection(config.Configs.DBConfig.UsersCollection)
	var result types.MongoUserDoc
	err := coll.FindOne(context.Background(), bson.D{bson.E{Key: "user_details.email", Value: email}}).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error finding the user: %v", err)
	}
	return result.UserDetails.PasswordHash, nil
}

func (s *StorageService) GetUserId(email string) (string, error) {
	coll := s.client.Database(config.Configs.DBConfig.DBName).Collection(config.Configs.DBConfig.UsersCollection)
	var result types.MongoUserDoc
	err := coll.FindOne(context.Background(), bson.D{bson.E{Key: "user_details.email", Value: email}}).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error finding the user: %v", err)
	}
	fmt.Println(result)
	return result.UserId, nil
}

func (s *StorageService) InsertUserDoc(userDoc *types.MongoUserDoc) error {
	coll := s.client.Database(config.Configs.DBConfig.DBName).Collection(config.Configs.DBConfig.UsersCollection)
	_, err := coll.InsertOne(context.Background(), userDoc)
	if err != nil {
		return fmt.Errorf("error inserting into db: %v", err)
	}
	return nil
}

func (s *StorageService) InsertAclDoc(aclDoc *types.MongoAclsDoc) error {
	coll := s.client.Database(config.Configs.DBConfig.DBName).Collection(config.Configs.DBConfig.AclsCollection)
	_, err := coll.InsertOne(context.Background(), aclDoc)
	if err != nil {
		return fmt.Errorf("error inserting into db: %v", err)
	}
	return nil
}
