package models

import "encoding/json"

func (req *CreateOrganizationRequest) ToDomain() *Organization {
	return &Organization{
		Name:        req.Name,
		Geom:        Point{Latitude: req.Latitude, Longitude: req.Longitude},
		Address:     req.Address,
		Description: req.Description,
		ImageUrls:   req.ImageUrls,
	}
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
