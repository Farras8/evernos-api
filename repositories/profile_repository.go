package repositories

import (
	"evernos-api2/models"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{
		db: db,
	}
}

func (r *profileRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *profileRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}