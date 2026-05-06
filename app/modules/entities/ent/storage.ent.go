package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Storage struct {
	bun.BaseModel `bun:"table:storages,alias:st"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Provider  string    `bun:"provider,notnull"`
	Path      string    `bun:"path,nullzero"`
	URL       string    `bun:"url,nullzero"`
	FileSize  int64     `bun:"file_size,notnull"`
	MimeType  string    `bun:"mime_type,nullzero"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	ShortCode string    `bun:"short_code,notnull"`
}
