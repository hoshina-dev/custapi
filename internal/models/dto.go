package models

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse is the DTO for user responses
type UserResponse struct {
	ID             string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email          string `json:"email" example:"user@example.com"`
	Name           string `json:"name" example:"John Doe"`
	OrganizationID string `json:"organization_id" example:"550e8400-e29b-41d4-a716-446655440001"`
}

// CreateUserRequest is the DTO for user creation
type CreateUserRequest struct {
	Email          string `json:"email" validate:"required,email" example:"user@example.com"`
	Name           string `json:"name" validate:"required" example:"John Doe"`
	OrganizationID string `json:"organization_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440001"`
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
	Address     *string  `json:"address,omitempty" validate:"omitempty" example:"254 St, Bangkok, TH"`
	Description *string  `json:"description,omitempty" validate:"omitempty" example:"Higher education institution"`
	ImageUrls   []string `json:"image_urls,omitempty" validate:"omitempty,dive,url" example:"https://example.com/example-1.jpg,https://example.com/example-2.jpg"`
}

type UpdateOrganizationRequest struct {
	Name        *string  `json:"name" validate:"omitempty" example:"Acme Corp"`
	Latitude    *float64 `json:"lat" validate:"omitempty,latitude" example:"13.7388"`
	Longitude   *float64 `json:"lng" validate:"omitempty,longitude" example:"100.5322"`
	Address     *string  `json:"address,omitempty" validate:"omitempty" example:"254 St, Bangkok, TH"`
	Description *string  `json:"description,omitempty" validate:"omitempty" example:"Higher education institution"`
	ImageUrls   []string `json:"image_urls,omitempty" validate:"omitempty,dive,url" example:"https://example.com/example-1.jpg,https://example.com/example-2.jpg"`
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
