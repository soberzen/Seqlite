package service

import (
	"context"
	"crypto/rand"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/repository"

	"github.com/google/uuid"
)

type CollaborationService struct {
	documents *repository.DocumentAccessPostgresRepository
	tickets   *repository.CollaborationTicketRedisRepository
}

type CollaborationTicket struct {
	Ticket  string `json:"ticket"`
	CanEdit bool   `json:"canEdit"`
}

func NewCollaborationService(
	documents *repository.DocumentAccessPostgresRepository,
	tickets *repository.CollaborationTicketRedisRepository,
) *CollaborationService {
	return &CollaborationService{
		documents: documents,
		tickets:   tickets,
	}
}

func (s *CollaborationService) requireRead(
	ctx context.Context,
	publicID uuid.UUID,
	userID uint,
) (*repository.DocumentAccess, error) {
	access, err := s.documents.HasDocumentPermission(
		ctx,
		publicID.String(),
		userID,
	)
	if err != nil {
		return nil, err
	}
	if !access.CanRead {
		return nil, errs.New(errs.KindForbidden, "document not readable")
	}
	return access, nil
}

func (s *CollaborationService) IssueTicket(
	ctx context.Context,
	publicID uuid.UUID,
	userID uint,
) (*CollaborationTicket, error) {
	access, err := s.requireRead(ctx, publicID, userID)
	if err != nil {
		return nil, err
	}
	ticket := rand.Text()
	err = s.tickets.Save(ctx, publicID, ticket, userID)
	if err != nil {
		return nil, err
	}
	return &CollaborationTicket{
		Ticket:  ticket,
		CanEdit: access.CanEdit,
	}, nil
}

type CollaborationAccess struct {
	DocumentID int64
	ActorID    int64
	CanEdit    bool
}

func (s *CollaborationService) AuthorizeConnect(
	ctx context.Context,
	publicID uuid.UUID,
	ticket string,
) (*CollaborationAccess, error) {
	userID, err := s.tickets.Consume(ctx, publicID, ticket)
	if err != nil {
		return nil, err
	}
	access, err := s.requireRead(ctx, publicID, userID)
	if err != nil {
		return nil, err
	}
	return &CollaborationAccess{
		DocumentID: access.DocumentID,
		ActorID:    int64(userID),
		CanEdit:    access.CanEdit,
	}, nil
}
