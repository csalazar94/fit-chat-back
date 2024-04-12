package service

import (
	"context"

	"github.com/csalazar94/fit-chat-back/internal/db"
)

type IUserService interface {
	Create(context.Context, CreateUserParams) (User, error)
	GetAll(context.Context) ([]User, error)
}

type IMessageService interface {
	Create(context.Context, CreateMessageParams) (Message, error)
}

type IAuthService interface {
	Login(context.Context, string, string) (bool, error)
}

type Services struct {
	UserService    IUserService
	AuthService    IAuthService
	MessageService IMessageService
}

func NewServices(dbQueries *db.Queries) *Services {
	return &Services{
		UserService:    NewUserService(dbQueries),
		AuthService:    NewAuthService(dbQueries),
		MessageService: NewMessageService(dbQueries),
	}
}
