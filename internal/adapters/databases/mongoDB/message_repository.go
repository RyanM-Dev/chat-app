package mongoDB

import (
	"chat-app/internal/core/domain"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

//type MessageRepository interface {
//	SendMessage(chatID, userID domain.ID, message string) error
//	DeleteMessage(chatID, userID, messageID domain.ID) error
//}

func (mr *MongoRepository) SendMessage(chatID, userID domain.ID, message string) error {
	coll := mr.GetCollection("messages")
	timeNow := time.Now()
	newMessage := domain.Message{
		SenderID:    userID,
		ChatID:      chatID,
		CreatedTime: &timeNow,
		Content:     message,
	}

	newMongoMessage, err := ToMongoMessage(newMessage)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(context.Background(), newMongoMessage)
	if err != nil {
		return err
	}

	return nil
}

func (mr *MongoRepository) DeleteMessage(messageID domain.ID) error {
	messageColl := mr.GetCollection("messages")
	filter := bson.M{"_id": messageID}
	delRes, err := messageColl.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if delRes.DeletedCount == 0 {
		return fmt.Errorf("message not found")
	}
	return nil
}

func (mr *MongoRepository) GetMessage(messageID domain.ID) (domain.Message, error) {
	messageColl := mr.GetCollection("messages")
	filter := bson.M{"_id": messageID}
	var mongoMessage MongoMessage
	err := messageColl.FindOne(context.Background(), filter).Decode(&mongoMessage)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Message{}, fmt.Errorf("message not found")
		}
		return domain.Message{}, err
	}
	message := mongoMessage.ToDomainMessage()

	return message, nil

}
