package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/csalazar94/fit-chat-back/internal/db"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
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
	dbQueries    *db.Queries
	openaiClient *openai.Client
}

func NewMessageService(dbQueries *db.Queries, openaiClient *openai.Client) *messageService {
	return &messageService{dbQueries, openaiClient}
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

const (
	SystemRoleId    = 1
	AssistantRoleId = 2
	ToolRoleId      = 3
	UserRoleId      = 4
)

func (messageService *messageService) AIMessageStream(ctx context.Context, chatId uuid.UUID) (*openai.ChatCompletionStream, error) {
	roleMap := map[int32]string{
		SystemRoleId:    "system",
		AssistantRoleId: "assistant",
		ToolRoleId:      "tool",
		UserRoleId:      "user",
	}

	dbMessages, err := messageService.dbQueries.GetMessagesByChatId(ctx, chatId)
	if err != nil {
		return nil, fmt.Errorf("error al obtener mensajes por id de chat: %v", err)
	}
	var messages []openai.ChatCompletionMessage
	for _, dbMessage := range dbMessages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    roleMap[dbMessage.AuthorRoleID],
			Content: dbMessage.Content,
		})
	}
	req := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo0125,
		Messages: messages,
		Stream:   true,
	}
	stream, err := messageService.openaiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return stream, fmt.Errorf("error al crear stream de mensajes: %v", err)
	}
	return stream, err
}
