package mongoDB

import (
	"chat-app/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MongoChat struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Owner       string             `bson:"owner"`
	Admins      []string           `bson:"admins"`
	Members     []string           `bson:"members"`
	CreatedTime *time.Time         `bson:"created_time"`
	DeletedTime *time.Time         `bson:"deleted_time"`
	ChatType    int                `bson:"chat_type"`
}

type MongoMessage struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	SenderID string             `bson:"sender_id"`
	ChatID   string             `bson:"chat_id"`
	Content  string             `bson:"content"`
}

func ToMongoChat(chat domain.Chat) MongoChat {
	admins := make([]string, len(chat.Admins))
	for i, admin := range chat.Admins {
		admins[i] = string(admin)
	}

	members := make([]string, len(chat.Members))
	for i, member := range chat.Members {
		members[i] = string(member)
	}

	return MongoChat{

		Name:        chat.Name,
		Owner:       string(chat.Owner),
		Admins:      admins,
		Members:     members,
		CreatedTime: chat.CreatedTime,
		DeletedTime: chat.DeletedTime,
		ChatType:    int(chat.ChatType),
	}
}

func (mc *MongoChat) ToDomainChat() domain.Chat {
	admins := make([]domain.ID, len(mc.Admins))
	for i, admin := range mc.Admins {
		admins[i] = domain.ID(admin)
	}

	members := make([]domain.ID, len(mc.Members))
	for i, member := range mc.Members {
		members[i] = domain.ID(member)
	}

	return domain.Chat{
		ID:          domain.ID(mc.ID.Hex()),
		Name:        mc.Name,
		Owner:       domain.ID(mc.Owner),
		Admins:      admins,
		Members:     members,
		CreatedTime: mc.CreatedTime,
		DeletedTime: mc.DeletedTime,
		ChatType:    domain.ChatType(mc.ChatType),
	}
}

func ToMongoMessage(message domain.Message) MongoMessage {
	return MongoMessage{
		SenderID: string(message.SenderID),
		ChatID:   string(message.ChatID),
		Content:  message.Content,
	}
}

func (mm *MongoMessage) ToDomainMessage() domain.Message {
	return domain.Message{
		ID:       domain.ID(mm.ID.Hex()),
		SenderID: domain.ID(mm.SenderID),
		ChatID:   domain.ID(mm.ChatID),
		Content:  mm.Content,
	}
}
