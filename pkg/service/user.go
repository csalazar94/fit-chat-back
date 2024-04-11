package service

import (
	"context"
	"log"
	"time"

	"github.com/csalazar94/fit-chat-back/internal/db"
	"github.com/csalazar94/fit-chat-back/pkg/password"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserParams struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type userService struct {
	dbQueries *db.Queries
}

func NewUserService(dbQueries *db.Queries) *userService {
	return &userService{dbQueries}
}

func (userService *userService) Create(context context.Context, params CreateUserParams) (user User, err error) {
	encodedHash, err := password.GenerateHash(params.Password, nil, nil)
	if err != nil {
		log.Printf("Error al generar hash: %v", err)
		return user, err
	}
	dbUser, err := userService.dbQueries.CreateUser(context, db.CreateUserParams{
		ID:          params.ID,
		FullName:    params.FullName,
		Email:       params.Email,
		EncodedHash: encodedHash,
		CreatedAt:   params.CreatedAt,
		UpdatedAt:   params.UpdatedAt,
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

func (userService *userService) GetAll(context context.Context) (users []User, err error) {
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
