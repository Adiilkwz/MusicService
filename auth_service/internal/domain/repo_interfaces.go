package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) (int64, error)
	Update(ctx context.Context, user *User) error
	SetResetCode(ctx context.Context, email, code string) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, limit, offset int) ([]*User, error)
}
