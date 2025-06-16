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

	ErrGetMessages = errors.New("failed to get all message between user")
	ErrGetMessageByID = errors.New("failed to get message")
)

type messageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) module.MessageService {
	return &messageService{
		repo: repo,
	}
}

func (s *messageService) Save(ctx context.Context, mReq *model.MessageReq) error {
	msg := mReq.ToMessage()
	if err := s.repo.AddMesssage(ctx, msg); err != nil {
		return ErrSaveMessage
	}

	log.Printf("[Service]: message {%s} created success\n", msg.ID)
	return nil
}

func (s *messageService) Update(ctx context.Context, mReq *model.MessageReq, id string) error {
	msg := mReq.ToMessage()
	if err := s.repo.UpdateMessage(ctx, msg, id); err != nil {
		return ErrUpdateMessage
	}

	log.Printf("[Service]: message {%s} updated success\n", msg.ID)
	return nil
}

func (s *messageService) Delete(ctx context.Context, id string) error {
	if err := s.repo.DeleteMessage(ctx, id, &model.Message{}); err != nil {
		return ErrDeleteMessage
	}

	log.Printf("[Service]: message {%s} deleted success\n", id)
	return nil
}

func (s *messageService) GetMessagesBetweenUser(ctx context.Context, senderID string, receiverID string) ([]model.MessageResp, error) {
	messages, err := s.repo.GetMessagesBetweenUser(ctx, senderID, receiverID)
	if err != nil {
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
	if err := s.repo.GetMessageByID(ctx, id, &msg); err != nil {
		return nil, ErrGetMessageByID
	}

	mResp := msg.ToMessageResp()
	log.Printf("[Service]: get message {%s} success\n", msg.ID)
	return mResp, nil
}
