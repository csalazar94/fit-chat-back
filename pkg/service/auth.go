package service

import (
	"context"
	"log"

	"github.com/csalazar94/fit-chat-back/internal/db"
	"github.com/csalazar94/fit-chat-back/pkg/password"
)

type authService struct {
	dbQueries *db.Queries
}

func NewAuthService(dbQueries *db.Queries) *authService {
	return &authService{dbQueries}
}

func (authService *authService) Login(context context.Context, email, pass string) (bool, error) {
	user, err := authService.dbQueries.GetUserByEmail(context, email)
	if err != nil {
		log.Printf("Error al obtener usuario por email: %v", err)
		return false, err
	}
	ok, err := password.CompareHash(user.EncodedHash, pass)
	if err != nil {
		log.Printf("Error al comparar hash: %v", err)
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}
