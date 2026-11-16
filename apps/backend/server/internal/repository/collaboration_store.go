package repository

import (
	"context"
	"database/sql"
	"errors"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/model"
	"qingmo-server/internal/realtime"
	"strconv"

	"github.com/Deln0r/ygo/persist"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DocumentStore struct {
	db *gorm.DB
}

func NewDocumentStore(db *gorm.DB) *DocumentStore {
	return &DocumentStore{db: db}
}

func parseDocName(docName string) (uint, error) {
	id, err := strconv.ParseUint(docName, 10, min(strconv.IntSize, 63))
	if err != nil {
		return 0, errs.New(errs.KindInvalid, "invalid document name")
	}
	return uint(id), nil
}

func (s *DocumentStore) withLockedDocument(
	ctx context.Context,
	docName string,
	run func(*gorm.DB, *model.Documents) error,
) error {
	id, err := parseDocName(docName)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var doc model.Documents
		err := tx.Select("id", "yjs_revision").
			Clauses(clause.Locking{Strength: "UPDATE"}).Take(&doc, id).Error
		if err != nil {
			return err
		}
		return run(tx, &doc)
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
}

func readDocumentUpdates(
	tx *gorm.DB,
	docID uint,
) ([][]byte, uint, uint, error) {
	var checkpoint model.DocumentCheckpoint

	err := tx.Where("document_id = ?", docID).
		Order("revision DESC").
		Take(&checkpoint).
		Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, 0, errs.Wrap(
			errs.KindInternal,
			"failed to query document checkpoint",
			err,
		)
	}
	var updates [][]byte
	last := checkpoint.Revision
	if err == nil {
		updates = append(updates, checkpoint.YjsState)
	}
	var revisions []model.DocumentRevision
	err = tx.Select("revision", "yjs_update").
		Where("document_id = ? AND revision > ?", docID, last).
		Order("revision ASC").
		Find(&revisions).Error
	if err != nil {
		return nil, 0, 0, errs.Wrap(
			errs.KindInternal,
			"failed to query document revisions",
			err,
		)
	}
	for _, rev := range revisions {
		updates = append(updates, rev.YjsUpdate)
		last = rev.Revision
	}
	return updates, last, uint(len(revisions)), nil
}

func (s *DocumentStore) StoreUpdate(
	ctx context.Context,
	docName string,
	update []byte,
) error {
	if len(update) == 0 {
		return persist.ErrEmptyUpdate
	}
	actorID, ok := realtime.GetActorID(ctx)
	if !ok {
		return errs.New(
			errs.KindUnauthorized,
			"missing authenticated update actor",
		)
	}
	return s.withLockedDocument(
		ctx,
		docName,
		func(tx *gorm.DB, doc *model.Documents) error {
			next := doc.YjsRevision + 1
			revision := model.DocumentRevision{
				DocumentID: doc.ID,
				Revision:   next,
				YjsUpdate:  update,
				ActorID:    actorID,
			}
			if err := tx.Create(&revision).Error; err != nil {
				return errs.Wrap(
					errs.KindInternal,
					"create document revision failed",
					err,
				)
			}
			if err := tx.Model(doc).
				UpdateColumn("yjs_revision", next).
				Error; err != nil {
				return errs.Wrap(
					errs.KindInternal,
					"update document revision failed",
					err,
				)
			}
			return nil
		},
	)

}

func (s *DocumentStore) GetUpdates(
	ctx context.Context,
	docName string,
) ([][]byte, error) {
	id, err := parseDocName(docName)
	if err != nil {
		return nil, err
	}
	var updates [][]byte
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var readErr error
		updates, _, _, readErr = readDocumentUpdates(tx, id)
		return readErr
	}, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	return updates, err

}
func (s *DocumentStore) Flush(ctx context.Context, docName string) error {
	err := s.withLockedDocument(
		ctx,
		docName,
		func(tx *gorm.DB, doc *model.Documents) error {
			updates, last, newRevisionCount, err := readDocumentUpdates(
				tx,
				doc.ID,
			)
			if err != nil {
				return errs.Wrap(
					errs.KindInternal,
					"no updates to flush",
					err,
				)
			}
			if newRevisionCount <= 2 {
				return errs.New(
					errs.KindInternal,
					"new revisions <= 3 can not flush",
				)
			}
			state, err := persist.MergeUpdates(updates)
			if err != nil {
				return errs.Wrap(
					errs.KindInternal,
					"failed to merge updates",
					err,
				)
			}
			if err := tx.Create(&model.DocumentCheckpoint{
				DocumentID: doc.ID,
				Revision:   last,
				YjsState:   state,
			}).Error; err != nil {
				return errs.Wrap(
					errs.KindInternal,
					"create document checkpoint failed",
					err,
				)
			}
			return nil
		},
	)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return errs.Wrap(
			errs.KindInternal,
			"failed to flush document",
			err,
		)
	}
	return nil
}

func (s *DocumentStore) DocumentExists(
	ctx context.Context,
	docName string,
) (bool, error) {
	id, err := parseDocName(docName)
	if err != nil {
		return false, err
	}
	var revision model.DocumentRevision
	err = s.db.WithContext(ctx).
		Select("document_id").
		Where("document_id = ?", id).
		Take(&revision).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errs.Wrap(
			errs.KindInternal,
			"failed to query document revision",
			err,
		)
	}
	return true, nil
}

func (s *DocumentStore) ListDocuments(ctx context.Context) ([]string, error) {
	var ids []uint

	err := s.db.WithContext(ctx).
		Model(&model.DocumentRevision{}).
		Distinct("document_id").
		Pluck("document_id", &ids).
		Error
	if err != nil {
		return nil, errs.Wrap(
			errs.KindInternal,
			"failed to query document ids",
			err,
		)
	}

	rooms := make([]string, 0, len(ids))
	for _, id := range ids {
		rooms = append(rooms, strconv.FormatUint(uint64(id), 10))
	}
	return rooms, nil
}

func (s *DocumentStore) ClearDocument(
	ctx context.Context,
	docName string,
) error {
	err := s.withLockedDocument(
		ctx,
		docName,
		func(tx *gorm.DB, doc *model.Documents) error {
			if err := tx.Where("document_id = ?", doc.ID).
				Delete(&model.DocumentCheckpoint{}).
				Error; err != nil {
				return errs.Wrap(
					errs.KindInternal,
					"failed to delete document checkpoints",
					err,
				)
			}
			if err := tx.Where("document_id = ?", doc.ID).
				Delete(&model.DocumentRevision{}).
				Error; err != nil {
				return errs.Wrap(
					errs.KindInternal,
					"failed to delete document revisions",
					err,
				)
			}
			return nil
		},
	)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return errs.Wrap(
			errs.KindInternal,
			"failed to clear document",
			err,
		)
	}
	return nil
}

func (s *DocumentStore) Close() error {
	return nil
}
