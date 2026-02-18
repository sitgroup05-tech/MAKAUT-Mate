package user

import (
	"context"
	"errors"

	"github.com/sitgroup05-tech/MakautMate/backend/internal/model"
	"gorm.io/gorm"
)

type sessionService struct {
	db *gorm.DB
}

func (s sessionService) GetUser(ctx context.Context, email string, u *model.User) error {

	if err := s.db.WithContext(ctx).First(u, "email = ?", email); err != nil {
		return errors.New("Invalid User Mail.")
	}

	return nil
}

func (s sessionService) ValidatePassword(opwd string, gpwd string) error {
	return nil
}
