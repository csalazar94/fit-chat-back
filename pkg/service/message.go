package service

import (
	"context"
	"log"
	"time"

	"github.com/csalazar94/fit-chat-back/internal/db"
	"github.com/google/uuid"
)

type Message struct {
	ID           uuid.UUID `json:"id"`
	ChatID       uuid.UUID `json:"chat_id"`
	AuthorRoleID int32     `json:"author_role_id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateMessageParams struct {
	ID           uuid.UUID `json:"id"`
	ChatID       uuid.UUID `json:"chat_id"`
	AuthorRoleID int32     `json:"author_role_id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type messageService struct {
	dbQueries *db.Queries
}

func NewMessageService(dbQueries *db.Queries) *messageService {
	return &messageService{dbQueries}
}

func (messageService *messageService) Create(context context.Context, params CreateMessageParams) (message Message, err error) {
	dbMessage, err := messageService.dbQueries.CreateMessage(context, db.CreateMessageParams{
		ID:           params.ID,
		ChatID:       params.ChatID,
		AuthorRoleID: params.AuthorRoleID,
		Content:      params.Content,
		CreatedAt:    params.CreatedAt,
		UpdatedAt:    params.UpdatedAt,
	})
	if err != nil {
		log.Printf("Error al crear mensaje: %v", err)
		return message, err
	}
	return Message{
		ID:           dbMessage.ID,
		ChatID:       dbMessage.ChatID,
		AuthorRoleID: dbMessage.AuthorRoleID,
		Content:      dbMessage.Content,
		CreatedAt:    dbMessage.CreatedAt,
		UpdatedAt:    dbMessage.UpdatedAt,
	}, nil
}
