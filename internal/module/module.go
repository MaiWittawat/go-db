package module

import (
	"context"
	"go-rebuild/internal/model"
)

type MessageService interface {
	Save(ctx context.Context, message *model.MessageReq) error
	Update(ctx context.Context, message *model.MessageReq, id string) error
	Delete(ctx context.Context, id string) error

	GetMessagesBetweenUser(ctx context.Context, senderID string, receiverID string) ([]model.MessageResp, error)
	GetMessageByID(ctx context.Context, id string) (*model.MessageResp, error)
}

type StockService interface {
	Save(ctx context.Context, productID string, quantity int) error
	Update(ctx context.Context, productID string, quantity int) error
	IncreaseQuantity(ctx context.Context, q int, id string) error
	DecreaseQuantity(ctx context.Context, q int, id string) error
	Delete(ctx context.Context, id string) error
}

type OrderService interface {
	Save(ctx context.Context, oReq *model.OrderReq, userID string) error
	Update(ctx context.Context, o *model.Order, id string) error
	Delete(ctx context.Context, id string, userID string) error

	GetAll(ctx context.Context) ([]model.OrderResp, error)
	GetByID(ctx context.Context, id string) (*model.OrderResp, error)
}

type ProductService interface {
	Save(ctx context.Context, p *model.ProductReq, userID string) error
	Update(ctx context.Context, p *model.ProductReq, id string, userID string) error
	Delete(ctx context.Context, id string) error

	GetAll(ctx context.Context) ([]model.ProductResp, error)
	GetByID(ctx context.Context, id string) (*model.ProductResp, error)
}

type UserService interface {
	Save(ctx context.Context, user *model.User) error
	Update(ctx context.Context, u *model.User, id string) error
	Delete(ctx context.Context, id string) error

	GetAll(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}
