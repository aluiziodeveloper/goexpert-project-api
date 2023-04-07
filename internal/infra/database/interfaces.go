package database

import "github.com/aluiziodeveloper/goexpert-project-api/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}
