package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

func (req *CreateOrganizationRequest) ToDomain() *Organization {
	return &Organization{
		Name:        req.Name,
		Geom:        Point{Latitude: req.Latitude, Longitude: req.Longitude},
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
		org.Geom.Latitude = req.Latitude
		org.Geom.Longitude = req.Longitude
	}
	org.Address = req.Address
	org.Description = req.Description
	return org
}

func (org *Organization) ToResponse() OrganizationResponse {
	return OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Geom:        json.RawMessage(org.Geom.Geom),
		Address:     org.Address,
		Description: org.Description,
		ImageUrls:   org.ImageUrls,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}
}
