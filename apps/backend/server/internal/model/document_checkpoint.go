package model

import "time"

// 这是文档的全量记录
type DocumentCheckpoint struct {
	DocumentID uint      `gorm:"column:document_id;type:bigint;primaryKey;autoIncrement:false;not null"`
	Revision   uint      `gorm:"column:revision;type:bigint;primaryKey;autoIncrement:false;not null"`
	YjsState   []byte    `gorm:"column:yjs_state;type:bytea;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime"`
}
