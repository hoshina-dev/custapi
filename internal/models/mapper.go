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

func (req *CreateUserRequest) ToDomain() (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		Email:              req.Email,
		Name:               req.Name,
		Password:           string(hashedPassword),
		PhoneNumber:        req.PhoneNumber,
		SocialMedia:        req.SocialMedia,
		Description:        req.Description,
		AvatarURL:          req.AvatarURL,
		ResearchCategories: req.ResearchCategories,
	}
	return user, nil
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
	return user, nil
}

func (user *User) ToResponse() UserResponse {
	orgs := make([]UserMembershipResponse, 0, len(user.Organizations))
	for _, m := range user.Organizations {
		orgs = append(orgs, UserMembershipResponse{
			OrganizationID:   m.OrganizationID,
			OrganizationName: m.Organization.Name,
			IsAdmin:          m.IsAdmin,
		})
	}
	return UserResponse{
		ID:                 user.ID,
		Email:              user.Email,
		Name:               user.Name,
		PhoneNumber:        user.PhoneNumber,
		SocialMedia:        user.SocialMedia,
		Description:        user.Description,
		AvatarURL:          user.AvatarURL,
		ResearchCategories: user.ResearchCategories,
		Organizations:      orgs,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
	}
}

func (user *User) ToDetailResponse() UserDetailResponse {
	return UserDetailResponse{
		ID:                 user.ID,
		Email:              user.Email,
		Name:               user.Name,
		OrganizationID:     user.OrganizationID,
		Password:           user.Password,
		IsAdmin:            user.IsAdmin,
		PhoneNumber:        user.PhoneNumber,
		SocialMedia:        user.SocialMedia,
		Description:        user.Description,
		AvatarURL:          user.AvatarURL,
		ResearchCategories: user.ResearchCategories,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
	}
}
