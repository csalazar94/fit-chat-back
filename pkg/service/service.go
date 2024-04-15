package service

import (
	"context"

	"github.com/csalazar94/fit-chat-back/internal/db"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type IUserService interface {
	Create(context.Context, CreateUserParams) (User, error)
	GetAll(context.Context) ([]User, error)
}

type IMessageService interface {
	Create(context.Context, CreateMessageParams) (Message, error)
	AIMessageStream(context.Context, uuid.UUID) (*openai.ChatCompletionStream, error)
	GetAllByChatID(context.Context, GetAllByChatIDParams) ([]Message, error)
}

type IAuthService interface {
	Login(context.Context, string, string) (bool, error)
}

type IChatService interface {
	Create(context.Context, CreateChatParams) (Chat, error)
	GetAll(context.Context) ([]Chat, error)
}

type Services struct {
	UserService    IUserService
	AuthService    IAuthService
	MessageService IMessageService
	ChatService    IChatService
}

func NewServices(dbQueries *db.Queries, openaiClient *openai.Client) *Services {
	return &Services{
		UserService:    NewUserService(dbQueries),
		AuthService:    NewAuthService(dbQueries),
		MessageService: NewMessageService(dbQueries, openaiClient),
		ChatService:    NewChatService(dbQueries),
	}
}
