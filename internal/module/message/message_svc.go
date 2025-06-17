package messsage

import (
	"context"
	"errors"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"go-rebuild/internal/repository"

	log "github.com/sirupsen/logrus"
)

var (
	ErrSaveMessage   = errors.New("failed to save message")
	ErrUpdateMessage = errors.New("failed to update message")
	ErrDeleteMessage = errors.New("failed to delete message")

	ErrGetMessages    = errors.New("failed to get all message between user")
	ErrGetMessageByID = errors.New("failed to get message")
)

type messageService struct {
	repo repository.MessageRepository
}

// ------------------------ Constructor ------------------------
func NewMessageService(repo repository.MessageRepository) module.MessageService {
	return &messageService{
		repo: repo,
	}
}

// ------------------------ Method Basic CUD ------------------------
func (s *messageService) Save(ctx context.Context, mReq *model.MessageReq) error {
	msg := mReq.ToMessage()
	var baseLogFields = log.Fields{
		"message_id": msg.ID,
		"layer":      "message_service",
		"method":     "message_save",
	}

	if err := s.repo.AddMesssage(ctx, msg); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("add message")
		return ErrSaveMessage
	}

	log.Printf("[Service]: message {%s} created success\n", msg.ID)
	return nil
}

func (s *messageService) Update(ctx context.Context, mReq *model.MessageReq, id string) error {
	msg := mReq.ToMessage()
	var baseLogFields = log.Fields{
		"message_id": msg.ID,
		"layer":      "message_service",
		"method":     "message_update",
	}

	if err := s.repo.UpdateMessage(ctx, msg, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("update message")
		return ErrUpdateMessage
	}

	log.Printf("[Service]: message {%s} updated success\n", msg.ID)
	return nil
}

func (s *messageService) Delete(ctx context.Context, id string) error {
	var baseLogFields = log.Fields{
		"message_id": id,
		"layer":      "message_service",
		"method":     "message_delete",
	}
	if err := s.repo.DeleteMessage(ctx, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("delete message")
		return ErrDeleteMessage
	}

	log.Printf("[Service]: message {%s} deleted success\n", id)
	return nil
}

// ------------------------ Method Basic Query ------------------------
func (s *messageService) GetMessagesBetweenUser(ctx context.Context, senderID string, receiverID string) ([]model.MessageResp, error) {
	var baseLogFields = log.Fields{
		"sender_id":   senderID,
		"receiver_id": receiverID,
		"layer":       "message_service",
		"method":      "message_delete",
	}
	messages, err := s.repo.GetMessagesBetweenUser(ctx, senderID, receiverID)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get message between user")
		return nil, ErrGetMessages
	}

	var messagesResp []model.MessageResp
	for _, message := range messages {
		mResp := message.ToMessageResp()
		messagesResp = append(messagesResp, *mResp)
	}

	log.Info("[Service]: get all message success")
	return messagesResp, nil
}

func (s *messageService) GetMessageByID(ctx context.Context, id string) (*model.MessageResp, error) {
	var msg model.Message
	var baseLogFields = log.Fields{
		"message_id": id,
		"layer":      "message_service",
		"method":     "message_delete",
	}
	if err := s.repo.GetMessageByID(ctx, id, &msg); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get message by id")
		return nil, ErrGetMessageByID
	}

	mResp := msg.ToMessageResp()
	log.Printf("[Service]: get message {%s} success\n", msg.ID)
	return mResp, nil
}
