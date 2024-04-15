package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/csalazar94/fit-chat-back/pkg/service"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"nhooyr.io/websocket"
)

type wsHandler struct {
	clientsMux    sync.Mutex
	clients       map[*client]struct{}
	eventHandlers map[string]eventHandler
	services      *service.Services
}

type client struct {
	conn *websocket.Conn
}

func NewWsHandler(services *service.Services) *wsHandler {
	wsHandler := &wsHandler{
		clients:       make(map[*client]struct{}),
		eventHandlers: map[string]eventHandler{},
		services:      services,
	}
	wsHandler.setupEventHandlers()
	return wsHandler
}

func (h *wsHandler) subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("Conexión establecida")
	defer conn.CloseNow()

	client := &client{
		conn: conn,
	}
	h.addClient(client)
	defer h.removeClient(client)

	for {
		_, msg, err := conn.Read(r.Context())

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			log.Printf("Conexión cerrada")
			return
		} else if err != nil {
			log.Printf("%v", err)
			return
		}

		var event event
		err = json.Unmarshal(msg, &event)
		if err != nil {
			log.Printf("Error al decodificar mensaje: %v", err)
			continue
		}

		err = h.routeEvent(event, client, r.Context())
		if err != nil {
			log.Printf("Error al enrutar evento: %v", err)
			continue
		}
	}
}

func (h *wsHandler) addClient(client *client) {
	h.clientsMux.Lock()
	defer h.clientsMux.Unlock()
	log.Printf("Clientes conectados %d", len(h.clients)+1)
	h.clients[client] = struct{}{}
}

func (h *wsHandler) removeClient(client *client) {
	h.clientsMux.Lock()
	defer h.clientsMux.Unlock()
	log.Printf("Clientes conectados %d", len(h.clients)-1)
	delete(h.clients, client)
}

type event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	sendMessageType     = "send_message"
	messageReceivedType = "message_received"
	partialResponseType = "partial_response"
)

type sendMessagePayload struct {
	service.CreateMessageParams
}

type messageReceivedPayload struct {
	ID uuid.UUID `json:"id"`
}

type partialResponsePayload struct {
	ID             uuid.UUID `json:"id"`
	PartialContent string    `json:"partial_content"`
}

func (h *wsHandler) setupEventHandlers() {
	h.eventHandlers[sendMessageType] = h.sendMessageHandler
}

type eventHandler func(event, *client, context.Context) error

func (h *wsHandler) routeEvent(event event, c *client, ctx context.Context) error {
	handler, ok := h.eventHandlers[event.Type]
	if !ok {
		return fmt.Errorf("evento no manejado: %v", event)
	}
	err := handler(event, c, ctx)
	if err != nil {
		return fmt.Errorf("error al manejar evento: %v", err)
	}
	return nil
}

func (h *wsHandler) sendMessageHandler(event event, client *client, ctx context.Context) error {
	var sendMessagePayload sendMessagePayload
	err := json.Unmarshal(event.Payload, &sendMessagePayload)
	if err != nil {
		return fmt.Errorf("error al decodificar payload: %v", err)
	}
	msg, err := h.services.MessageService.Create(ctx, service.CreateMessageParams{
		ID:           sendMessagePayload.ID,
		ChatID:       sendMessagePayload.ChatID,
		AuthorRoleID: sendMessagePayload.AuthorRoleID,
		Content:      sendMessagePayload.Content,
		CreatedAt:    sendMessagePayload.CreatedAt,
		UpdatedAt:    sendMessagePayload.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("error al crear mensaje: %v", err)
	}
	err = h.sendMessageReceived(ctx, messageReceivedPayload{msg.ID}, client)
	if err != nil {
		return fmt.Errorf("error al enviar mensaje recibido: %v", err)
	}

	stream, err := h.services.MessageService.AIMessageStream(ctx, msg.ChatID)
	if err != nil {
		return fmt.Errorf("error al obtener stream de IA: %v", err)
	}
	defer stream.Close()
	responseId := uuid.New()
	var responseContent strings.Builder
	for {
		var response openai.ChatCompletionStreamResponse
		response, err = stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("error al recibir respuesta de IA: %v", err)
		}
		responseContent.WriteString(response.Choices[0].Delta.Content)
		err = h.sendPartialResponse(ctx, partialResponsePayload{responseId, responseContent.String()}, client)
		if err != nil {
			return fmt.Errorf("error al enviar respuesta parcial: %v", err)
		}
	}
	_, err = h.services.MessageService.Create(ctx, service.CreateMessageParams{
		ID:           responseId,
		ChatID:       sendMessagePayload.ChatID,
		AuthorRoleID: service.AssistantRoleId,
		Content:      responseContent.String(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("error al crear mensaje respuesta: %v", err)
	}
	return nil
}

func (h *wsHandler) sendEvent(ctx context.Context, event event, client *client) error {
	eventJson, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error al codificar evento: %v", err)
	}
	err = client.conn.Write(ctx, websocket.MessageText, []byte(eventJson))
	if err != nil {
		return fmt.Errorf("error al enviar mensaje: %v", err)
	}
	return nil
}

func (h *wsHandler) sendMessageReceived(ctx context.Context, messageReceivedPayload messageReceivedPayload, client *client) error {
	eventPayload, err := json.Marshal(messageReceivedPayload)
	if err != nil {
		return fmt.Errorf("error al codificar payload: %v", err)
	}
	event := event{
		Type:    messageReceivedType,
		Payload: eventPayload,
	}
	return h.sendEvent(ctx, event, client)
}

func (h *wsHandler) sendPartialResponse(ctx context.Context, partialResponsePayload partialResponsePayload, client *client) error {
	eventPayload, err := json.Marshal(partialResponsePayload)
	if err != nil {
		return fmt.Errorf("error al codificar payload: %v", err)
	}
	event := event{
		Type:    partialResponseType,
		Payload: eventPayload,
	}
	return h.sendEvent(ctx, event, client)
}
