CREATE TABLE IF NOT EXISTS storages (
    id UUID PRIMARY KEY,
    bucket_name VARCHAR(255) NOT NULL,
    object_key TEXT NOT NULL,
    file_name TEXT NOT NULL,
    content_type VARCHAR(255) NOT NULL,
    size_bytes BIGINT NOT NULL,
    size_text VARCHAR(64) NOT NULL,
    public_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS storages_object_key_uq ON storages (object_key);
CREATE INDEX IF NOT EXISTS storages_deleted_at_idx ON storages (deleted_at);
CREATE INDEX IF NOT EXISTS storages_created_at_idx ON storages (created_at DESC);

COMMENT ON TABLE storages IS 'Stores uploaded file metadata for MinIO objects and public access links.';
COMMENT ON COLUMN storages.id IS 'Storage record identifier (UUID).';
COMMENT ON COLUMN storages.bucket_name IS 'MinIO bucket where the object is stored.';
COMMENT ON COLUMN storages.object_key IS 'Object key/path inside the MinIO bucket.';
COMMENT ON COLUMN storages.file_name IS 'Original file name after basic sanitization.';
COMMENT ON COLUMN storages.content_type IS 'MIME type of the uploaded file.';
COMMENT ON COLUMN storages.size_bytes IS 'File size in bytes from uploaded object metadata.';
COMMENT ON COLUMN storages.size_text IS 'Human-readable file size (for example, 1.25 MB).';
COMMENT ON COLUMN storages.public_url IS 'Public URL used by internal clients to access object path.';
COMMENT ON COLUMN storages.created_at IS 'Record creation timestamp.';
COMMENT ON COLUMN storages.updated_at IS 'Record update timestamp.';
COMMENT ON COLUMN storages.deleted_at IS 'Soft delete timestamp; NULL means active.';
