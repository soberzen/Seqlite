package model

import "time"

// 这是文档的增量记录
type DocumentRevision struct {
	DocumentID uint      `gorm:"column:document_id;type:bigint;primaryKey;autoIncrement:false;not null"`
	Revision   uint      `gorm:"column:revision;type:bigint;primaryKey;autoIncrement:false;not null"`
	YjsUpdate  []byte    `gorm:"column:yjs_update;type:bytea;not null"`
	ActorID    uint      `gorm:"column:actor_id;type:bigint;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime"`
}

func (DocumentRevision) TableName() string { return "document_revisions" }
