package repo

import "gorm.io/gorm"

type Repository interface {
}

var _Repository = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &UserRepository{
		db: db,
	}
}
