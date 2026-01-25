package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (req *CreateOrganizationRequest) ToDomain() *Organization {
	return &Organization{
		Name:        req.Name,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Address:     req.Address,
		Description: req.Description,
		ImageUrls:   req.ImageUrls,
	}
}

func (req *UpdateOrganizationRequest) ToDomain(id uuid.UUID) *Organization {
	org := &Organization{ID: id, Address: req.Address,
		Description: req.Description, ImageUrls: req.ImageUrls}
	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Latitude != nil && req.Longitude != nil {
		org.Latitude = req.Latitude
		org.Longitude = req.Longitude
	}
	org.Address = req.Address
	org.Description = req.Description
	return org
}

func (org *Organization) ToResponse() OrganizationResponse {
	return OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Latitude:    *org.Latitude,
		Longitude:   *org.Longitude,
		Address:     org.Address,
		Description: org.Description,
		ImageUrls:   org.ImageUrls,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}
}

func (req *UpdateUserRequest) ToDomain(id uuid.UUID) (*User, error) {
	user := &User{ID: id, PhoneNumber: req.PhoneNumber, SocialMedia: req.SocialMedia, Description: req.Description,
		AvatarURL: req.AvatarURL, ResearchCategories: req.ResearchCategories}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.OrganizationID != nil {
		user.OrganizationID = *req.OrganizationID
	}
	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}
	return user, nil
}
