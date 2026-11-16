package service

import (
	"context"

	"qingmo-server/internal/dto"
	"qingmo-server/internal/model"
	"qingmo-server/internal/repository"

	"github.com/google/uuid"
)

type SpaceService struct {
	spaces *repository.SpacePostgresRepository
}

func NewSpaceService(spaces *repository.SpacePostgresRepository) *SpaceService {
	return &SpaceService{spaces: spaces}
}

func (s *SpaceService) Create(
	ctx context.Context,
	userID uint,
	req dto.CreateSpaceRequest,
) (*dto.SpaceResponse, error) {
	space, err := s.spaces.Create(ctx, &model.Spaces{
		Name:        req.Name,
		Emoji:       req.Emoji,
		Description: req.Description,
		OwnerID:     userID,
	})
	if err != nil {
		return nil, err
	}
	return spaceResponse(space), nil
}

func (s *SpaceService) ListVisibleSpaces(
	ctx context.Context,
	userID uint,
) ([]dto.ListSpacesResponse, error) {
	spaces, err := s.spaces.ListVisible(ctx, userID)
	if err != nil {
		return nil, err
	}
	return spaceListResponse(spaces), nil
}

func (s *SpaceService) ListTrashedSpaces(
	ctx context.Context,
	userID uint,
) ([]dto.ListSpacesResponse, error) {
	spaces, err := s.spaces.ListTrashed(ctx, userID)
	if err != nil {
		return nil, err
	}
	return spaceListResponse(spaces), nil
}

func (s *SpaceService) GetSpace(
	ctx context.Context,
	userID uint,
	publicID uuid.UUID,
) (*dto.SpaceDetailResponse, error) {
	space, err := s.spaces.GetVisibleByPublicID(ctx, userID, publicID)
	if err != nil {
		return nil, err
	}
	return &dto.SpaceDetailResponse{
		SpaceResponse: *spaceResponse(&space.Space),
		Role:          string(space.Role),
	}, nil
}

func (s *SpaceService) UpdateSpaceInfo(
	ctx context.Context,
	userID uint,
	publicID uuid.UUID,
	req dto.UpdateSpaceRequest,
) (*dto.SpaceResponse, error) {
	space, err := s.spaces.UpdateSpace(ctx, userID, publicID, req)
	if err != nil {
		return nil, err
	}
	return spaceResponse(space), nil
}

func (s *SpaceService) ListDocuments(
	ctx context.Context,
	userID uint,
	publicID uuid.UUID,
	req dto.ListSpaceDocumentsRequest,
) (*dto.ListSpaceDocumentsResponse, error) {
	docs, hasMore, err := s.spaces.ListVisibleSpaceDocument(
		ctx, userID, publicID, req.Page, req.PageSize,
	)
	if err != nil {
		return nil, err
	}
	items := make([]dto.SpaceDocumentResponse, 0, len(docs))
	for _, doc := range docs {
		items = append(items, dto.SpaceDocumentResponse{
			ID:    doc.PublicID,
			Title: doc.Title,
			Emoji: doc.Emoji,
		})
	}
	return &dto.ListSpaceDocumentsResponse{
		Items:    items,
		Page:     req.Page,
		PageSize: req.PageSize,
		HasMore:  hasMore,
	}, nil
}

func spaceResponse(space *model.Spaces) *dto.SpaceResponse {
	return &dto.SpaceResponse{
		ID:          space.PublicID,
		Name:        space.Name,
		Emoji:       space.Emoji,
		Description: space.Description,
	}
}

func spaceListResponse(spaces []*model.Spaces) []dto.ListSpacesResponse {
	items := make([]dto.ListSpacesResponse, 0, len(spaces))
	for _, space := range spaces {
		items = append(items, dto.ListSpacesResponse{
			ID:    space.PublicID,
			Name:  space.Name,
			Emoji: space.Emoji,
		})
	}
	return items
}
