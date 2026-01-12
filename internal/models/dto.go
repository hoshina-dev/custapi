package models

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
	ID   string `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name string `json:"name" example:"Acme Corp"`
}

// CreateOrganizationRequest is the DTO for organization creation
type CreateOrganizationRequest struct {
	Name string `json:"name" validate:"required" example:"Acme Corp"`
}

// ErrorResponse is the DTO for error responses
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
