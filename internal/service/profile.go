package service

import (
	"strings"

	"github.com/google/uuid"

	"dotsat.work/internal/model"
	"dotsat.work/internal/repository"
	"dotsat.work/internal/validation"
)

type ProfileService struct {
	profileRepo repository.ProfileRepository
}

func NewProfileService(profileRepo repository.ProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
	}
}

func (s *ProfileService) ByUserID(userID uuid.UUID) (*model.Profile, error) {
	return s.profileRepo.ByUserID(userID)
}

func (s *ProfileService) Create(profile *model.Profile) error {
	return s.profileRepo.Create(profile)
}

func (s *ProfileService) UpdateName(userID uuid.UUID, name string) error {
	name = strings.TrimSpace(name)

	err := validation.ValidateName(name)
	if err != nil {
		return err
	}

	return s.profileRepo.UpdateName(userID, name)
}
