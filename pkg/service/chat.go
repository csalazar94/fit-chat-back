package service

import (
	"context"
	"fmt"
	"time"

	"github.com/csalazar94/fit-chat-back/internal/db"
	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type chatService struct {
	dbQueries *db.Queries
}

func NewChatService(dbQueries *db.Queries) *chatService {
	return &chatService{dbQueries}
}

func (chatService *chatService) GetAll(ctx context.Context) (chats []Chat, err error) {
	dbChats, err := chatService.dbQueries.GetChats(ctx)
	if err != nil {
		return chats, fmt.Errorf("error al obtener chats: %w", err)
	}
	for _, dbChat := range dbChats {
		chats = append(chats, Chat{
			ID:        dbChat.ID,
			UserID:    dbChat.UserID,
			Title:     dbChat.Title.String,
			CreatedAt: dbChat.CreatedAt,
			UpdatedAt: dbChat.UpdatedAt,
		})
	}
	return chats, nil
}
