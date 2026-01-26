package models

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse is the DTO for user responses
type UserResponse struct {
	ID                 uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email              string    `json:"email" example:"user@example.com"`
	Name               string    `json:"name" example:"John Doe"`
	OrganizationID     uuid.UUID `json:"organization_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	IsAdmin            bool      `json:"is_admin" example:"true"`
	PhoneNumber        *string   `json:"phone_number,omitempty" example:"+1234567890"`
	SocialMedia        *string   `json:"social_media,omitempty" example:"@john on Twitter, linkedin.com/in/john"`
	Description        *string   `json:"description,omitempty" example:"Senior researcher specializing in quantum computing"`
	AvatarURL          *string   `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
	ResearchCategories []string  `json:"research_categories" example:"QuantumComputing,Qiskit,Cryogenics"`
	CreatedAt          time.Time `json:"created_at" example:"2026-01-01T12:00:00.00000+07:00"`
	UpdatedAt          time.Time `json:"updated_at" example:"2026-01-01T12:00:00.00000+07:00"`
}

// CreateUserRequest is the DTO for user creation
type CreateUserRequest struct {
	Email              string    `json:"email" validate:"required,email" example:"user@example.com"`
	Name               string    `json:"name" validate:"required,max=255" example:"John Doe"`
	OrganizationID     uuid.UUID `json:"organization_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440001"`
	Password           string    `json:"password" validate:"required" example:"PassWord123!"`
	PhoneNumber        *string   `json:"phone_number" validate:"omitempty,e164" example:"+1234567890"`
	SocialMedia        *string   `json:"social_media" example:"@john on Twitter, linkedin.com/in/john"`
	Description        *string   `json:"description" example:"Senior researcher specializing in quantum computing"`
	AvatarURL          *string   `json:"avatar_url" validate:"omitempty,url" example:"https://example.com/avatar.jpg"`
	ResearchCategories []string  `json:"research_categories" example:"QuantumComputing,Qiskit,Cryogenics"`
	IsAdmin            *bool     `json:"is_admin" example:"true"`
}

type UpdateUserRequest struct {
	Email              *string    `json:"email" validate:"omitempty,email" example:"user@example.com"`
	Name               *string    `json:"name" validate:"omitempty,max=255" example:"John Doe"`
	OrganizationID     *uuid.UUID `json:"organization_id" validate:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440001"`
	Password           *string    `json:"password" example:"PassWord123!"`
	PhoneNumber        *string    `json:"phone_number" validate:"omitempty,e164" example:"+1234567890"`
	SocialMedia        *string    `json:"social_media" example:"@john on Twitter, linkedin.com/in/john"`
	Description        *string    `json:"description" example:"Senior researcher specializing in quantum computing"`
	AvatarURL          *string    `json:"avatar_url" validate:"omitempty,url" example:"https://example.com/avatar.jpg"`
	ResearchCategories []string   `json:"research_categories" example:"QuantumComputing,Qiskit,Cryogenics"`
	IsAdmin            *bool      `json:"is_admin" example:"true"`
}

// OrganizationResponse is the DTO for organization responses
type OrganizationResponse struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name        string    `json:"name" example:"Acme Corp"`
	Latitude    float64   `json:"lat" example:"13.7888"`
	Longitude   float64   `json:"lng" example:"100.5322"`
	Address     *string   `json:"address,omitempty" example:"254 St, Bangkok, TH"`
	Description *string   `json:"description,omitempty" example:"Higher education institution"`
	ImageUrls   []string  `json:"image_urls" example:"https://example.com/example-1.jpg,https://example.com/example-2.jpg"`
	CreatedAt   time.Time `json:"created_at" example:"2026-01-01T12:00:00.00000+07:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2026-01-01T12:00:00.00000+07:00"`
}

// CreateOrganizationRequest is the DTO for organization creation
type CreateOrganizationRequest struct {
	Name        string   `json:"name" validate:"required" example:"Acme Corp"`
	Latitude    *float64 `json:"lat" validate:"required,latitude" example:"13.7388"`
	Longitude   *float64 `json:"lng" validate:"required,longitude" example:"100.5322"`
	Address     *string  `json:"address" example:"254 St, Bangkok, TH"`
	Description *string  `json:"description" example:"Higher education institution"`
	ImageUrls   []string `json:"image_urls" validate:"omitempty,dive,url" example:"https://example.com/example-1.jpg,https://example.com/example-2.jpg"`
}

type UpdateOrganizationRequest struct {
	Name        *string  `json:"name" example:"Acme Corp"`
	Latitude    *float64 `json:"lat" validate:"omitempty,latitude" example:"13.7388"`
	Longitude   *float64 `json:"lng" validate:"omitempty,longitude" example:"100.5322"`
	Address     *string  `json:"address" example:"254 St, Bangkok, TH"`
	Description *string  `json:"description" example:"Higher education institution"`
	ImageUrls   []string `json:"image_urls" validate:"omitempty,dive,url" example:"https://example.com/example-1.jpg,https://example.com/example-2.jpg"`
}

type OrganizationCoord struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Latitude  float64   `json:"lat" example:"13.7388"`
	Longitude float64   `json:"lng" example:"100.5322"`
}

type GetOrganizationsByIDsRequest struct {
	IDs []uuid.UUID `json:"ids" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000,550e8400-e29b-41d4-a716-446655440001"`
}

// ErrorResponse is the DTO for error responses
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
