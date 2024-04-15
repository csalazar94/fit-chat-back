package service

import (
	"context"
	"database/sql"
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

type CreateChatParams struct {
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

func (chatService *chatService) Create(ctx context.Context, params CreateChatParams) (chat Chat, err error) {
	var title string
	var titleIsValid bool
	if params.Title != "" {
		title = params.Title
		titleIsValid = true
	}
	dbChat, err := chatService.dbQueries.CreateChat(ctx, db.CreateChatParams{
		ID:        params.ID,
		UserID:    params.UserID,
		Title:     sql.NullString{String: title, Valid: titleIsValid},
		CreatedAt: params.CreatedAt,
		UpdatedAt: params.UpdatedAt,
	})
	if err != nil {
		return chat, fmt.Errorf("error al crear chat: %w", err)
	}
	return Chat{
		ID:        dbChat.ID,
		UserID:    dbChat.UserID,
		Title:     dbChat.Title.String,
		CreatedAt: dbChat.CreatedAt,
		UpdatedAt: dbChat.UpdatedAt,
	}, nil
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
