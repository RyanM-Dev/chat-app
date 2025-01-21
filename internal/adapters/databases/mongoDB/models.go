package mongoDB

import (
	"chat-app/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MongoChat struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Name        string               `bson:"name"`
	Owner       primitive.ObjectID   `bson:"owner"`
	Admins      []primitive.ObjectID `bson:"admins"`
	Members     []primitive.ObjectID `bson:"members"`
	CreatedTime *time.Time           `bson:"created_time"`
	DeletedTime *time.Time           `bson:"deleted_time"`
	ChatType    int                  `bson:"chat_type"`
}

type MongoMessage struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	SenderID primitive.ObjectID `bson:"sender_id"`
	ChatID   primitive.ObjectID `bson:"chat_id"`
	Content  string             `bson:"content"`
}

type MongoUser struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Username    string               `bson:"username"`
	FirstName   string               `bson:"first_name"`
	LastName    string               `bson:"last_name"`
	Password    string               `bson:"password"`
	Gender      int                  `bson:"gender"`
	Email       string               `bson:"email"`
	Contacts    []primitive.ObjectID `bson:"contacts"`
	DateOfBirth *time.Time           `bson:"date_of_birth"`
	CreatedTime *time.Time           `bson:"created_time"`
	DeletedTime *time.Time           `bson:"deleted_time"`
}

func ToMongoChat(chat domain.Chat) MongoChat {
	admins := make([]primitive.ObjectID, len(chat.Admins))
	for i, admin := range chat.Admins {
		mongoAdmin, err := primitive.ObjectIDFromHex(string(admin))
		if err != nil {
			panic(err)
		}
		admins[i] = mongoAdmin
	}

	members := make([]primitive.ObjectID, len(chat.Members))
	for i, member := range chat.Members {
		mongoMember, err := primitive.ObjectIDFromHex(string(member))
		if err != nil {
			panic(err)
		}
		members[i] = mongoMember
	}

	mongoOwner, err := primitive.ObjectIDFromHex(string(chat.Owner))
	if err != nil {
		panic(err)
	}

	return MongoChat{
		Name:        chat.Name,
		Owner:       mongoOwner,
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
		admins[i] = domain.ID(admin.Hex())
	}

	members := make([]domain.ID, len(mc.Members))
	for i, member := range mc.Members {
		members[i] = domain.ID(member.Hex())
	}

	return domain.Chat{
		ID:          domain.ID(mc.ID.Hex()),
		Name:        mc.Name,
		Owner:       domain.ID(mc.Owner.Hex()),
		Admins:      admins,
		Members:     members,
		CreatedTime: mc.CreatedTime,
		DeletedTime: mc.DeletedTime,
		ChatType:    domain.ChatType(mc.ChatType),
	}
}

func ToMongoMessage(message domain.Message) MongoMessage {
	senderID, err := primitive.ObjectIDFromHex(string(message.SenderID))
	if err != nil {
		panic(err)
	}

	chatID, err := primitive.ObjectIDFromHex(string(message.ChatID))
	if err != nil {
		panic(err)
	}

	return MongoMessage{
		SenderID: senderID,
		ChatID:   chatID,
		Content:  message.Content,
	}
}

func (mm *MongoMessage) ToDomainMessage() domain.Message {
	return domain.Message{
		ID:       domain.ID(mm.ID.Hex()),
		SenderID: domain.ID(mm.SenderID.Hex()),
		ChatID:   domain.ID(mm.ChatID.Hex()),
		Content:  mm.Content,
	}
}

func ToMongoUser(user domain.User) MongoUser {
	contacts := make([]primitive.ObjectID, len(user.Contacts))
	for i, contact := range user.Contacts {
		mongoContact, err := primitive.ObjectIDFromHex(string(contact))
		if err != nil {
			panic(err)
		}
		contacts[i] = mongoContact
	}

	return MongoUser{
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Password:    user.Password,
		Gender:      int(user.Gender),
		Email:       user.Email,
		Contacts:    contacts,
		DateOfBirth: user.DateOfBirth,
		CreatedTime: user.CreatedTime,
		DeletedTime: user.DeletedTime,
	}
}

func (mu *MongoUser) ToDomainUser() domain.User {
	contacts := make([]domain.ID, len(mu.Contacts))
	for i, contact := range mu.Contacts {
		contacts[i] = domain.ID(contact.Hex())
	}

	return domain.User{
		ID:          domain.ID(mu.ID.Hex()),
		Username:    mu.Username,
		FirstName:   mu.FirstName,
		LastName:    mu.LastName,
		Password:    mu.Password,
		Gender:      domain.Gender(mu.Gender),
		Email:       mu.Email,
		Contacts:    contacts,
		DateOfBirth: mu.DateOfBirth,
		CreatedTime: mu.CreatedTime,
		DeletedTime: mu.DeletedTime,
	}
}
