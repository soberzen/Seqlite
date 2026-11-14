package repository

import (
	"context"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentAccess struct {
	DocumentID int64 `gorm:"column:document_id"`
	CanRead    bool  `gorm:"column:can_read"`
	CanEdit    bool  `gorm:"column:can_edit"`
}

type DocumentAccessPostgresRepository struct {
	db *gorm.DB
}

func NewDocumentAccessPostgresRepository(
	db *gorm.DB,
) *DocumentAccessPostgresRepository {
	return &DocumentAccessPostgresRepository{
		db: db,
	}
}

func (d *DocumentAccessPostgresRepository) HasDocumentPermission(
	ctx context.Context,
	publicID string,
	userID uint,
) (*DocumentAccess, error) {
	id, err := uuid.Parse(publicID)
	if err != nil || id == uuid.Nil {
		return nil, errs.New(errs.KindInvalid, "invalid document id")
	}
	var rows []DocumentAccess

	result := d.db.WithContext(ctx).Raw(`
		SELECT doc.id AS document_id,
				doc.owner_id = u.id OR COALESCE(bool_or(p.code = 'document:read'), false) AS can_read,
				doc.owner_id = u.id OR COALESCE(bool_or(p.code = 'document:update'), false) AS can_edit
		FROM documents AS doc
		JOIN users AS u ON u.id = ? AND u.deleted_at IS NULL
		LEFT JOIN document_grants AS g ON g.document_id = doc.id AND g.user_id = u.id
		LEFT JOIN role_document_permissions AS rp ON rp.role_id = g.role_id
		LEFT JOIN document_permissions AS p ON p.id = rp.permission_id
		WHERE doc.public_id = ? AND doc.status = ?
		GROUP BY doc.id, doc.owner_id, u.id;
	`, userID, id, model.DocumentStatusNormal).Scan(&rows)

	if result.Error != nil {
		return nil, errs.Wrap(
			errs.KindInternal,
			"failed to query document access",
			result.Error,
		)
	}
	if len(rows) == 0 {
		return nil, errs.New(errs.KindNotFound, "document not found")
	}
	return &rows[0], nil
}
