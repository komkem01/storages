ALTER TABLE storages
    ADD COLUMN IF NOT EXISTS bucket_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS object_key TEXT,
    ADD COLUMN IF NOT EXISTS file_name TEXT,
    ADD COLUMN IF NOT EXISTS content_type VARCHAR(255),
    ADD COLUMN IF NOT EXISTS size_bytes BIGINT,
    ADD COLUMN IF NOT EXISTS size_text VARCHAR(64),
    ADD COLUMN IF NOT EXISTS public_url TEXT,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

UPDATE storages
SET bucket_name = COALESCE(bucket_name, ''),
    object_key = COALESCE(object_key, path),
    file_name = COALESCE(file_name, path),
    content_type = COALESCE(content_type, mime_type),
    size_bytes = COALESCE(size_bytes, file_size),
    size_text = COALESCE(size_text, file_size::text || ' B'),
    public_url = COALESCE(public_url, url);

ALTER TABLE storages
    ALTER COLUMN bucket_name SET NOT NULL,
    ALTER COLUMN object_key SET NOT NULL,
    ALTER COLUMN file_name SET NOT NULL,
    ALTER COLUMN content_type SET NOT NULL,
    ALTER COLUMN size_bytes SET NOT NULL,
    ALTER COLUMN size_text SET NOT NULL,
    ALTER COLUMN public_url SET NOT NULL;

DROP INDEX IF EXISTS storages_short_code_uq;
DROP INDEX IF EXISTS storages_provider_idx;

ALTER TABLE storages
    DROP COLUMN IF EXISTS provider,
    DROP COLUMN IF EXISTS path,
    DROP COLUMN IF EXISTS url,
    DROP COLUMN IF EXISTS file_size,
    DROP COLUMN IF EXISTS mime_type,
    DROP COLUMN IF EXISTS short_code;

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
