package dto

import "github.com/google/uuid"

type CreateSpaceRequest struct {
	Name        string  `json:"name"        binding:"required,min=6,max=255,displayname"`
	Emoji       *string `json:"emoji"       binding:"omitempty,emoji"`
	Description string  `json:"description" binding:"omitempty,max=280"`
}

type SpaceResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Emoji       *string   `json:"emoji"`
	Description string    `json:"description"`
}

type SpaceDetailResponse struct {
	SpaceResponse
	Role string `json:"role"`
}

type ListSpacesResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Emoji *string   `json:"emoji"`
}

type UpdateSpaceRequest struct {
	Name        *string `json:"name"        binding:"omitempty,min=6,max=255,displayname"`
	Emoji       *string `json:"emoji"       binding:"omitempty,emoji"`
	Description *string `json:"description" binding:"omitempty,max=280"`
}

type ListSpaceDocumentsRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"pageSize,default=20"`
}

type SpaceDocumentResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Emoji *string   `json:"emoji"`
}

type ListSpaceDocumentsResponse struct {
	Items    []SpaceDocumentResponse `json:"items"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"pageSize"`
	HasMore  bool                    `json:"hasMore"`
}
