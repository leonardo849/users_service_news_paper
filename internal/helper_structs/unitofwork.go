package helper_structs

import (
	"users-service/internal/repository"
	"users-service/internal/unitofwork"

	"gorm.io/gorm"
)

func CreateUnitOfWork(db *gorm.DB, userStatusRepository *repository.UserStatusRepository, userRepository *repository.UserRepository) *unitofwork.Unitofwork {
	return unitofwork.CreateUnitOfWork(userRepository, userStatusRepository, db)
}