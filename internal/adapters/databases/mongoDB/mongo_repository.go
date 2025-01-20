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
	"time"
)

func (mr *MongoRepository) CreateChat(chat domain.Chat) (domain.ID, error) {
	coll := mr.GetCollection("chats")

	mongoChat := ToMongoChat(chat)

	res, err := coll.InsertOne(context.Background(), mongoChat)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", fmt.Errorf("error creating chat in collection 'chats': %v : %w", err, repositories.ErrDuplicateChat)
		}
		return "", err
	}
	chatID := res.InsertedID.(primitive.ObjectID).Hex()

	return domain.ID(chatID), nil

}

func (mr *MongoRepository) FindChat(chatID domain.ID) (domain.Chat, error) {
	coll := mr.GetCollection("chats")
	filter := bson.M{"_id": chatID}
	var mongoChat MongoChat
	err := coll.FindOne(context.Background(), filter).Decode(&mongoChat)
	if err != nil {
		return domain.Chat{}, err
	}

	chat := mongoChat.ToDomainChat()
	return chat, nil

}

func (mr *MongoRepository) UpdateChatName(chat domain.Chat) error {
	coll := mr.GetCollection("chats")
	filter := bson.M{"_id": chat.ID}
	update := bson.M{"$set": bson.M{"name": chat.Name}}
	res, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating chat name for chat %v : %w", chat.ID, err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("error finding chat ID %v to update:%w ", chat.ID, repositories.ErrChatNotFound)
	}
	return nil
}

func (mr *MongoRepository) DeleteChat(chatID domain.ID) error {
	coll := mr.GetCollection("chats")
	filter := bson.M{"_id": chatID}
	res, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error deleting chat %v : %v", chatID, err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("error deleting chat %v : %w ", chatID, repositories.ErrChatNotFound)
	}
	return nil
}

func (mr *MongoRepository) GetMessages(chatID domain.ID) ([]domain.Message, error) {
	messageCollection := mr.GetCollection("messages")
	cursor, err := messageCollection.Find(context.Background(), bson.M{"chat_id": chatID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var messages []domain.Message
	var mongoMessages []MongoMessage

	if err := cursor.All(context.Background(), &mongoMessages); err != nil {
		return nil, err
	}
	for _, mongoMessage := range mongoMessages {
		message := mongoMessage.ToDomainMessage()
		messages = append(messages, message)
	}
	return messages, nil
}

func (mr *MongoRepository) AddUser(chatID domain.ID, userIDs []domain.ID) error {

}

func (mr *MongoRepository) RemoveUser(chatID domain.ID, userID []domain.ID) error {

}

func (mr *MongoRepository) GetMembers(chatID domain.ID) ([]domain.ID, error) {

}

func (mr *MongoRepository) SetAdmin(userID, chatID domain.ID) error {

}
