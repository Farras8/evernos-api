package services

import (
	"evernos-api2/models"
	"evernos-api2/repositories"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ProfileService interface {
	GetProfile(userID uint) (*models.User, error)
	UpdateProfile(userID uint, updateData map[string]string) (*models.User, error)
}

type profileService struct {
	profileRepo repositories.ProfileRepository
}

func NewProfileService(profileRepo repositories.ProfileRepository) ProfileService {
	return &profileService{
		profileRepo: profileRepo,
	}
}

func (s *profileService) GetProfile(userID uint) (*models.User, error) {
	user, err := s.profileRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *profileService) UpdateProfile(userID uint, updateData map[string]string) (*models.User, error) {
	// Ambil user yang akan diupdate
	user, err := s.profileRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update field yang diberikan
	if updateData["nama"] != "" {
		user.Nama = updateData["nama"]
	}
	if updateData["email"] != "" {
		user.Email = updateData["email"]
	}
	if updateData["noTelp"] != "" {
		user.NoTelp = updateData["noTelp"]
	}
	if updateData["tanggalLahir"] != "" {
		tanggalLahir, err := time.Parse("2006-01-02", updateData["tanggalLahir"])
		if err != nil {
			return nil, errors.New("invalid date format for tanggalLahir. Use YYYY-MM-DD format")
		}
		user.TanggalLahir = tanggalLahir
	}
	if updateData["jenisKelamin"] != "" {
		user.JenisKelamin = updateData["jenisKelamin"]
	}
	if updateData["tentang"] != "" {
		user.Tentang = updateData["tentang"]
	}
	if updateData["pekerjaan"] != "" {
		user.Pekerjaan = updateData["pekerjaan"]
	}
	if updateData["idProvinsi"] != "" {
		user.IdProvinsi = updateData["idProvinsi"]
	}
	if updateData["idKota"] != "" {
		user.IdKota = updateData["idKota"]
	}

	// Update password jika diberikan
	if updateData["password"] != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData["password"]), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.KataSandi = string(hashedPassword)
	}

	// Simpan perubahan
	err = s.profileRepo.Update(user)
	if err != nil {
		return nil, errors.New("failed to update profile")
	}

	return user, nil
}