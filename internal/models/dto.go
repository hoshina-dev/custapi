package models

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse is the DTO for user responses
type UserResponse struct {
	ID                 uuid.UUID `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email              string    `json:"email" validate:"required" example:"user@example.com"`
	Name               string    `json:"name" validate:"required" example:"John Doe"`
	Role               UserRole  `json:"role" validate:"required" example:"user"`
	PhoneNumber        *string   `json:"phone_number,omitempty" example:"+1234567890"`
	SocialMedia        *string   `json:"social_media,omitempty" example:"@john on Twitter, linkedin.com/in/john"`
	Description        *string   `json:"description,omitempty" example:"Senior researcher specializing in quantum computing"`
	AvatarURL          *string   `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
	ResearchCategories []string  `json:"research_categories" validate:"required" example:"QuantumComputing,Qiskit,Cryogenics"`
	CreatedAt          time.Time `json:"created_at" validate:"required" example:"2026-01-01T12:00:00.00000+07:00"`
	UpdatedAt          time.Time `json:"updated_at" validate:"required" example:"2026-01-01T12:00:00.00000+07:00"`
} //	@name	UserResponse

// UserDetailResponse is the DTO for single user responses
type UserDetailResponse struct {
	ID                 uuid.UUID `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email              string    `json:"email" validate:"required" example:"user@example.com"`
	Name               string    `json:"name" validate:"required" example:"John Doe"`
	Role               UserRole  `json:"role" validate:"required" example:"user"`
	PhoneNumber        *string   `json:"phone_number,omitempty" example:"+1234567890"`
	SocialMedia        *string   `json:"social_media,omitempty" example:"@john on Twitter, linkedin.com/in/john"`
	Description        *string   `json:"description,omitempty" example:"Senior researcher specializing in quantum computing"`
	AvatarURL          *string   `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
	ResearchCategories []string  `json:"research_categories" validate:"required" example:"QuantumComputing,Qiskit,Cryogenics"`
	CreatedAt          time.Time `json:"created_at" validate:"required" example:"2026-01-01T12:00:00.00000+07:00"`
	UpdatedAt          time.Time `json:"updated_at" validate:"required" example:"2026-01-01T12:00:00.00000+07:00"`
} //	@name	UserDetailResponse

// VerifyCredentialsRequest is the DTO for verifying a user's email + password
type VerifyCredentialsRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"PassWord123!"`
} //	@name	VerifyCredentialsRequest

// UserWithRoleResponse is the DTO for user responses within an organization context, includes member role
type UserWithRoleResponse struct {
	UserResponse
	MemberRole MemberRole `json:"member_role" example:"user"`
} //	@name	UserWithRoleResponse

// UserMembershipResponse represents a user's membership in an organization
type UserMembershipResponse struct {
	OrganizationID uuid.UUID            `json:"organization_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Organization   OrganizationResponse `json:"organization"`
	Role           MemberRole           `json:"role" example:"user"`
	CreatedAt      time.Time            `json:"created_at" example:"2026-01-01T12:00:00.00000+07:00"`
} //	@name	UserMembershipResponse
type CreateUserRequest struct {
	Email              string   `json:"email" validate:"required,email" example:"user@example.com"`
	Name               string   `json:"name" validate:"required,max=255" example:"John Doe"`
	Password           string   `json:"password" validate:"required" example:"PassWord123!"`
	PhoneNumber        *string  `json:"phone_number" validate:"omitempty,e164" example:"+1234567890"`
	SocialMedia        *string  `json:"social_media" example:"@john on Twitter, linkedin.com/in/john"`
	Description        *string  `json:"description" example:"Senior researcher specializing in quantum computing"`
	AvatarURL          *string  `json:"avatar_url" validate:"omitempty,url" example:"https://example.com/avatar.jpg"`
	ResearchCategories []string `json:"research_categories" example:"QuantumComputing,Qiskit,Cryogenics"`
} //	@name	CreateUserRequest

type UpdateUserRequest struct {
	Email              *string  `json:"email" validate:"omitempty,email" example:"user@example.com"`
	Name               *string  `json:"name" validate:"omitempty,max=255" example:"John Doe"`
	Password           *string  `json:"password" example:"PassWord123!"`
	PhoneNumber        *string  `json:"phone_number" validate:"omitempty,e164" example:"+1234567890"`
	SocialMedia        *string  `json:"social_media" example:"@john on Twitter, linkedin.com/in/john"`
	Description        *string  `json:"description" example:"Senior researcher specializing in quantum computing"`
	AvatarURL          *string  `json:"avatar_url" validate:"omitempty,url" example:"https://example.com/avatar.jpg"`
	ResearchCategories []string `json:"research_categories" example:"QuantumComputing,Qiskit,Cryogenics"`
} //	@name	UpdateUserRequest

// AddMemberRequest is the DTO for adding a user to an organization
type AddMemberRequest struct {
	UserID uuid.UUID  `json:"user_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role   MemberRole `json:"role" example:"user"`
} //	@name	AddMemberRequest

// SetRoleRequest is the DTO for updating a member's role within an organization
type SetRoleRequest struct {
	Role MemberRole `json:"role" validate:"required" example:"manager"`
} //	@name	SetRoleRequest

// OrganizationResponse is the DTO for organization responses
type OrganizationResponse struct {
	ID          uuid.UUID `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name        string    `json:"name" validate:"required" example:"Acme Corp"`
	Latitude    float64   `json:"lat" validate:"required" example:"13.7888"`
	Longitude   float64   `json:"lng" validate:"required" example:"100.5322"`
	Address     *string   `json:"address,omitempty" example:"254 St, Bangkok, TH"`
	Description *string   `json:"description,omitempty" example:"Higher education institution"`
	ImageUrls   []string  `json:"image_urls" validate:"required" example:"https://example.com/example-1.jpg,https://example.com/example-2.jpg"`
	CreatedAt   time.Time `json:"created_at" validate:"required" example:"2026-01-01T12:00:00.00000+07:00"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required" example:"2026-01-01T12:00:00.00000+07:00"`
} //	@name	OrganizationResponse

// CreateOrganizationRequest is the DTO for organization creation
type CreateOrganizationRequest struct {
	Name        string   `json:"name" validate:"required" example:"Acme Corp"`
	Latitude    *float64 `json:"lat" validate:"required,latitude" example:"13.7388"`
	Longitude   *float64 `json:"lng" validate:"required,longitude" example:"100.5322"`
	Address     *string  `json:"address" example:"254 St, Bangkok, TH"`
	Description *string  `json:"description" example:"Higher education institution"`
	ImageUrls   []string `json:"image_urls" validate:"omitempty,dive,url" example:"https://example.com/example-1.jpg,https://example.com/example-2.jpg"`
} //	@name	CreateOrganizationRequest

type UpdateOrganizationRequest struct {
	Name        *string  `json:"name" example:"Acme Corp"`
	Latitude    *float64 `json:"lat" validate:"omitempty,latitude" example:"13.7388"`
	Longitude   *float64 `json:"lng" validate:"omitempty,longitude" example:"100.5322"`
	Address     *string  `json:"address" example:"254 St, Bangkok, TH"`
	Description *string  `json:"description" example:"Higher education institution"`
	ImageUrls   []string `json:"image_urls" validate:"omitempty,dive,url" example:"https://example.com/example-1.jpg,https://example.com/example-2.jpg"`
} //	@name	UpdateOrganizationRequest

type OrganizationCoord struct {
	ID        uuid.UUID `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
	Latitude  float64   `json:"lat" validate:"required" example:"13.7388"`
	Longitude float64   `json:"lng" validate:"required" example:"100.5322"`
} //	@name	OrganizationCoord

type GetOrganizationsByIDsRequest struct {
	IDs []uuid.UUID `json:"ids" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000,550e8400-e29b-41d4-a716-446655440001"`
} //	@name	GetOrganizationsByIDsRequest

// ErrorResponse is the DTO for error responses
type ErrorResponse struct {
	Error string `json:"error" validate:"required" example:"error message"`
} //	@name	ErrorResponse
