package mongoDB

import (
	"chat-app/internal/core/domain"
	"chat-app/internal/core/repositories"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (mr *MongoRepository) Register(user domain.User) (domain.ID, error) {
	coll := mr.GetCollection("users")
	mongoUser := ToMongoUser(user)
	res, err := coll.InsertOne(context.Background(), mongoUser)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %v : %w", err, repositories.ErrOperationFailed)
	}
	userID := domain.ID(res.InsertedID.(primitive.ObjectID).Hex())

	return userID, nil
}

func (mr *MongoRepository) Login(username, password string) (domain.ID, error) {
	coll := mr.GetCollection("users")
	filter := bson.M{"username": username, "password": password}

	projection := bson.M{"chatID": 1}

	var result struct {
		ChatID string `bson:"chatID"`
	}

	err := coll.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", fmt.Errorf("username or password incorrect: %v : %w", err, repositories.ErrWrongLoginInfo)
	}

	return domain.ID(result.ChatID), nil
}

func (mr *MongoRepository) GetChatIDList(userID domain.ID) ([]domain.ID, error) {
	coll := mr.GetCollection("users")

	mongoUserID, err := primitive.ObjectIDFromHex(string(userID))
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": mongoUserID}
	projection := bson.M{"chat_id_list": 1}

	var result struct {
		ChatIDList []primitive.ObjectID `bson:"chat_id_list"`
	}

	err = coll.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get chat ID list: %w", err)
	}

	// Convert MongoDB ObjectIDs to domain IDs
	var domainChatIDList []domain.ID
	for _, chatID := range result.ChatIDList {
		domainChatIDList = append(domainChatIDList, domain.ID(chatID.Hex()))
	}
	return domainChatIDList, nil
}

func (mr *MongoRepository) AddContact(userID, contactID domain.ID) error {
	coll := mr.GetCollection("users")

	// Convert userID and contactID to MongoDB ObjectID
	mongoUserID, err := primitive.ObjectIDFromHex(string(userID))
	if err != nil {
		return fmt.Errorf("invalid user ID: %v : %w", err, repositories.ErrInvalidID)
	}
	mongoContactID, err := primitive.ObjectIDFromHex(string(contactID))
	if err != nil {
		return fmt.Errorf("invalid contact ID: %v : %w", err, repositories.ErrInvalidID)
	}

	// Update the user's contact list
	filter := bson.M{"_id": mongoUserID}
	update := bson.M{"$addToSet": bson.M{"contacts": mongoContactID}}
	res, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to add contact: %w", err)
	}

	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil

}

func (mr *MongoRepository) RemoveContact(userID, contactID domain.ID) error {
	coll := mr.GetCollection("users")

	// Convert userID and contactID to MongoDB ObjectID
	mongoUserID, err := primitive.ObjectIDFromHex(string(userID))
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	mongoContactID, err := primitive.ObjectIDFromHex(string(contactID))
	if err != nil {
		return fmt.Errorf("invalid contact ID: %w", err)
	}

	// Update the user's contact list
	filter := bson.M{"_id": mongoUserID}
	update := bson.M{"$pull": bson.M{"contacts": mongoContactID}}
	res, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove contact: %w", err)
	}

	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
