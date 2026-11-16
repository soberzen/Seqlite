package repository

import (
	"context"
	"errors"
	"math"
	"qingmo-server/internal/dto"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SpacePostgresRepository struct {
	db *gorm.DB
}

func NewSpacePostgresRepository(db *gorm.DB) *SpacePostgresRepository {
	return &SpacePostgresRepository{db: db}
}

func (d *SpacePostgresRepository) Create(
	ctx context.Context,
	space *model.Spaces,
) (*model.Spaces, error) {
	if err := d.db.WithContext(ctx).
		Model(&model.Spaces{}).
		Create(space).
		Error; err != nil {
		return nil, errs.Wrap(errs.KindInternal, "create space error", err)
	}
	return space, nil
}

func (d *SpacePostgresRepository) ListVisible(
	ctx context.Context,
	userID uint,
) ([]*model.Spaces, error) {
	db := d.db.WithContext(ctx)
	joinedSpaceIDs := db.Model(&model.SpaceGrants{}).
		Select("space_grants.space_id").
		Joins("JOIN role_space_permissions AS rp ON rp.role_id = space_grants.role_id").
		Joins("JOIN space_permissions AS p ON p.id = rp.permission_id").
		Where("space_grants.user_id = ? AND p.code = 'space:read'", userID)

	var spaces []*model.Spaces
	if err := db.Where("status = ?", model.SpaceStatusNormal).
		Where("(owner_id = ? OR id IN (?))", userID, joinedSpaceIDs).
		Find(&spaces).
		Error; err != nil {
		return nil, errs.Wrap(
			errs.KindInternal,
			"list visible spaces error",
			err,
		)
	}
	return spaces, nil
}

func (d *SpacePostgresRepository) ListTrashed(
	ctx context.Context,
	userID uint,
) ([]*model.Spaces, error) {
	var spaces []*model.Spaces
	if err := d.db.WithContext(ctx).
		Where("owner_id = ? AND status = ?", userID, model.SpaceStatusTrashed).
		Find(&spaces).
		Error; err != nil {
		return nil, errs.Wrap(
			errs.KindInternal,
			"list trashed spaces error",
			err,
		)
	}
	return spaces, nil
}

type spaceWithRole struct {
	Space   model.Spaces   `gorm:"embedded"`
	Role    model.RoleCode `gorm:"column:role"`
	CanRead bool           `gorm:"column:can_read"`
}

func (d *SpacePostgresRepository) GetVisibleByPublicID(
	ctx context.Context,
	userID uint,
	publicID uuid.UUID,
) (*spaceWithRole, error) {
	var result spaceWithRole

	query := d.db.WithContext(ctx).Raw(`
		SELECT s.*,
			CASE WHEN s.owner_id = ? THEN ? ELSE COALESCE(r.code, '') END AS role,
			(s.owner_id = ? OR EXISTS (
				SELECT 1
				FROM role_space_permissions AS rp
				JOIN space_permissions AS p ON p.id = rp.permission_id
				WHERE rp.role_id = r.id AND p.code = 'space:read'
			)) AS can_read
		FROM spaces AS s
		LEFT JOIN space_grants AS g ON g.space_id = s.id AND g.user_id = ?
		LEFT JOIN roles AS r ON r.id = g.role_id
		WHERE s.public_id = ? AND s.status = ?
	`, userID, model.RoleOwner, userID, userID, publicID, model.SpaceStatusNormal).
		Scan(&result)
	if query.Error != nil {
		return nil, errs.Wrap(
			errs.KindInternal,
			"query space failed",
			query.Error,
		)
	}
	if query.RowsAffected == 0 {
		return nil, errs.New(errs.KindNotFound, "space not found")
	}

	if !result.CanRead {
		return nil, errs.New(
			errs.KindForbidden,
			"space read permission required",
		)
	}

	return &result, nil
}

func (d *SpacePostgresRepository) UpdateSpace(
	ctx context.Context,
	userID uint,
	publicID uuid.UUID,
	update dto.UpdateSpaceRequest,
) (*model.Spaces, error) {
	values := make(map[string]any, 3)
	if update.Name != nil {
		values["name"] = *update.Name
	}
	if update.Emoji != nil {
		values["emoji"] = *update.Emoji
	}
	if update.Description != nil {
		values["description"] = *update.Description
	}

	var space model.Spaces
	query := d.db.WithContext(ctx).Model(&space).
		Where("public_id = ? AND status = ?", publicID, model.SpaceStatusNormal).
		Where(`(owner_id = ? OR EXISTS (
			SELECT 1
			FROM space_grants AS g
			JOIN role_space_permissions AS rp ON rp.role_id = g.role_id
			JOIN space_permissions AS p ON p.id = rp.permission_id
			WHERE g.space_id = spaces.id AND g.user_id = ? AND p.code = 'space:update'
		))`, userID, userID)

	var result *gorm.DB
	if len(values) == 0 {
		result = query.Find(&space)
	} else {
		result = query.Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "public_id"}, {Name: "name"}, {Name: "emoji"},
			{Name: "description"},
		}}).Updates(values)
	}

	if err := result.Error; err != nil {
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, errs.Wrap(
			errs.KindInternal,
			"update space info failed",
			err,
		)
	}
	if result.RowsAffected == 0 {
		return nil, errs.New(errs.KindNotFound, "space not found")
	}
	return &space, nil
}

const maxPageSize = 100 // 单页上限，防止大分页拖垮 DB/内存
func (d *SpacePostgresRepository) ListVisibleSpaceDocument(
	ctx context.Context,
	userID uint,
	publicID uuid.UUID,
	page, pageSize int,
) ([]model.Documents, bool, error) {

	if page < 1 || pageSize < 1 || pageSize > maxPageSize ||
		page-1 > math.MaxInt/pageSize {
		return nil, false, errs.New(errs.KindInvalid, "invalid pagination")
	}

	docs := make([]model.Documents, 0, pageSize+1)

	err := d.db.WithContext(ctx).Raw(`
    SELECT doc.public_id, doc.title, doc.emoji
    FROM documents AS doc
    JOIN spaces AS s ON s.id = doc.space_id
    JOIN users AS u ON u.id = ? AND u.deleted_at IS NULL
    WHERE s.public_id = ?
      AND s.status = ?
      AND doc.status = ?
      AND (
        doc.owner_id = u.id
        OR EXISTS (
          SELECT 1
          FROM document_grants AS g
          JOIN role_document_permissions AS rp ON rp.role_id = g.role_id
          JOIN document_permissions AS p ON p.id = rp.permission_id
          WHERE g.document_id = doc.id
            AND g.user_id = u.id
            AND p.code = 'document:read'
        )
      )
    ORDER BY doc.id DESC
    LIMIT ? OFFSET ?
    `,
		userID, publicID, model.SpaceStatusNormal, model.DocumentStatusNormal,
		pageSize+1, (page-1)*pageSize,
	).Scan(&docs).Error

	if err != nil {
		return nil, false, errs.Wrap(
			errs.KindInternal,
			"list documents failed",
			err,
		)
	}
	hasMore := len(docs) > pageSize
	if hasMore {
		docs = docs[:pageSize]
	}
	return docs, hasMore, nil
}
