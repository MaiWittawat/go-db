package message

import (
	"context"
	"go-rebuild/internal/cache"
	dbRepo "go-rebuild/internal/db"
	"go-rebuild/internal/model"
	"go-rebuild/internal/repository"
)

type messageRepo struct {
	db         dbRepo.DB
	collection string
	cacheSvc   cache.Cache
	KeyGen     *cache.KeyGenerator
}

func NewMessageRepo(db dbRepo.DB, cacheSvc cache.Cache) repository.MessageRepository {
	keyGen := cache.NewKeyGenerator("messages")
	return &messageRepo{
		db:         db,
		collection: "messages",
		cacheSvc:   cacheSvc,
		KeyGen:     keyGen,
	}
}

func (r *messageRepo) AddMesssage(ctx context.Context, msg *model.Message) error {
	return r.db.Create(ctx, r.collection, msg)
}

func (r *messageRepo) UpdateMessage(ctx context.Context, msg *model.Message, id string) error {
	return r.db.Update(ctx, r.collection, msg, id)
}

func (r *messageRepo) DeleteMessage(ctx context.Context, id string, msg *model.Message) error {
	return r.db.Delete(ctx, r.collection, msg, id)
}

func (r *messageRepo) GetMessagesBetweenUser(ctx context.Context, senderID string, receiverID string) ([]model.Message, error) {
	messages, err := r.db.FindMessageBetweenUser(ctx, senderID, receiverID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *messageRepo) GetMessageByID(ctx context.Context, id string, msg *model.Message) error {
	if err := r.db.GetByID(ctx, r.collection, id, msg); err != nil {
		return err
	}
	return nil
}
