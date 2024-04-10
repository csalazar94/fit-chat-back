package service

import (
	"context"
	"log"
	"time"

	"github.com/csalazar94/fit-chat-back/internal/db"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserService struct {
	dbQueries *db.Queries
}

func NewUserService(dbQueries *db.Queries) *UserService {
	return &UserService{dbQueries}
}

func (userService *UserService) Create(context context.Context, params db.CreateUserParams) (user User, err error) {
	dbUser, err := userService.dbQueries.CreateUser(context, db.CreateUserParams{
		ID:        params.ID,
		FullName:  params.FullName,
		Email:     params.Email,
		Password:  params.Password,
		CreatedAt: params.CreatedAt,
		UpdatedAt: params.UpdatedAt,
	})
	if err != nil {
		log.Println("Error al crear usuario: ", err)
		return user, err
	}
	return User{
		ID:        dbUser.ID,
		FullName:  dbUser.FullName,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}, err
}

func (userService *UserService) GetAll(context context.Context) (users []User, err error) {
	dbUsers, err := userService.dbQueries.GetUsers(context)
	if err != nil {
		log.Println("Error al obtener usuarios: ", err)
		return nil, err
	}
	for _, dbUser := range dbUsers {
		user := User{
			ID:        dbUser.ID,
			FullName:  dbUser.FullName,
			Email:     dbUser.Email,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.CreatedAt,
		}
		users = append(users, user)
	}
	return users, err
}
