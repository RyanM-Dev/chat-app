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
	"log"
)

func (mr *MongoRepository) CreateChat(chat domain.Chat) (domain.ID, error) {
	chatColl := mr.GetCollection("chats")

	mongoChat := ToMongoChat(chat)

	res, err := chatColl.InsertOne(context.Background(), mongoChat)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", fmt.Errorf("error creating chat in collection 'chats': %v : %w", err, repositories.ErrDuplicateChat)
		}
		return "", err
	}
	chatID := res.InsertedID.(primitive.ObjectID).Hex()

	userColl := mr.GetCollection("users")
	mongoUserID := mongoChat.Owner
	filter := bson.M{"_id": mongoUserID}
	update := bson.M{"$addToSet": bson.M{"chatIDList": mongoChat.ID}}
	updateRes, err := userColl.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "", err
	}
	if updateRes.MatchedCount == 0 {
		return "", fmt.Errorf("failed to insert chat ID into user chat ID list")
	}

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

	userColl := mr.GetCollection("users")
	mongoChatID, err := primitive.ObjectIDFromHex(string(chatID))
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	userFilter := bson.M{"chatIDList": mongoChatID}
	userUpdate := bson.M{"$pull": bson.M{"chatIDList": mongoChatID}}
	updateRes, err := userColl.UpdateMany(context.Background(), userFilter, userUpdate)
	if err != nil {
		return fmt.Errorf("error removing chat from users list : %v", err)
	}
	if updateRes.MatchedCount == 0 {
		return fmt.Errorf("no users found with the speceified chat ID %v", chatID)
	}

	return nil
}

func (mr *MongoRepository) GetMessages(chatID domain.ID) ([]domain.Message, error) {
	messageCollection := mr.GetCollection("messages")
	cursor, err := messageCollection.Find(context.Background(), bson.M{"chat_id": chatID})
	if err != nil {
		return nil, err
	}
	defer func() {
		err := cursor.Close(context.Background())
		if err != nil {
			log.Printf("error closing cursor : %v", err)
		}

	}()

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
	chatColl := mr.GetCollection("chats")

	mongoChatID, err := primitive.ObjectIDFromHex(string(chatID))
	if err != nil {
		return err
	}

	mongoUserIDs := make([]primitive.ObjectID, len(userIDs))
	for _, userID := range userIDs {
		mongoUserID, err := primitive.ObjectIDFromHex(string(userID))
		if err != nil {
			return err
		}
		mongoUserIDs = append(mongoUserIDs, mongoUserID)
	}

	filter := bson.M{"_id": mongoChatID}

	res, err := chatColl.UpdateMany(context.Background(), filter, bson.M{"$addToSet": bson.M{"users": mongoUserIDs}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("error adding user to chat,no match was found")
	}

	userColl := mr.GetCollection("users")

	for _, mongoUserID := range mongoUserIDs {
		userFilter := bson.M{"_id": mongoUserID}
		userUpdate := bson.M{"$addToSet": bson.M{"chatIDList": mongoChatID}}
		updateRes, err := userColl.UpdateMany(context.Background(), userFilter, userUpdate)
		if err != nil {
			return err
		}
		if updateRes.MatchedCount == 0 {
			return fmt.Errorf("failed to insert chat ID into user chat ID list")
		}

	}

	return nil
}

func (mr *MongoRepository) RemoveUser(chatID domain.ID, userIDs []domain.ID) error {
	coll := mr.GetCollection("chats")

	mongoChatID, err := primitive.ObjectIDFromHex(string(chatID))
	if err != nil {
		return err
	}

	mongoUserIDs := make([]primitive.ObjectID, len(userIDs))
	for _, userID := range userIDs {
		mongoUserID, err := primitive.ObjectIDFromHex(string(userID))
		if err != nil {
			return err
		}
		mongoUserIDs = append(mongoUserIDs, mongoUserID)
	}

	filter := bson.M{"_id": mongoChatID}
	update := bson.M{"$pull": bson.M{"members": bson.M{"$in": mongoUserIDs}}}
	res, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("error removing user from chat,no match was found")
	}

	userColl := mr.GetCollection("users")

	for _, mongoUserID := range mongoUserIDs {
		userFilter := bson.M{"_id": mongoUserID}
		userUpdate := bson.M{"$pull": bson.M{"chatIDList": mongoChatID}}
		updateRes, err := userColl.UpdateMany(context.Background(), userFilter, userUpdate)
		if err != nil {
			return err
		}
		if updateRes.MatchedCount == 0 {
			return fmt.Errorf("failed to insert chat ID into user chat ID list")
		}

	}

	return nil
}

func (mr *MongoRepository) GetMembers(chatID domain.ID) ([]domain.ID, error) {
	coll := mr.GetCollection("chats")
	filter := bson.M{"_id": chatID}
	var result struct {
		Members []primitive.ObjectID `bson:"members"`
	}
	err := coll.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("failed to get members: chat with ID %s not found: %w", chatID, err)
		}
		return nil, fmt.Errorf("failed to get members from MongoDB: %w", err)
	}
	var members []domain.ID
	for _, member := range result.Members {
		var id domain.ID
		id = domain.ID(member.Hex())
		members = append(members, id)
	}
	return members, nil
}

func (mr *MongoRepository) SetAdmin(adminID, chatID domain.ID) error {
	coll := mr.GetCollection("chats")
	mongoChatID, err := primitive.ObjectIDFromHex(string(chatID))
	if err != nil {
		return err
	}
	mongoAdminID, err := primitive.ObjectIDFromHex(string(adminID))
	if err != nil {
		return err
	}

	filter := bson.M{"_id": mongoChatID}
	update := bson.M{"$addToSet": bson.M{"admins": mongoAdminID}}
	res, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("error adding admin to chat,no match was found")
	}
	return nil
}
