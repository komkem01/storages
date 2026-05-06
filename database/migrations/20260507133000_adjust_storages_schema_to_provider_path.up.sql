ALTER TABLE storages
    ADD COLUMN IF NOT EXISTS provider VARCHAR(100) DEFAULT 'Railway',
    ADD COLUMN IF NOT EXISTS path TEXT,
    ADD COLUMN IF NOT EXISTS url TEXT,
    ADD COLUMN IF NOT EXISTS file_size BIGINT,
    ADD COLUMN IF NOT EXISTS mime_type VARCHAR(255),
    ADD COLUMN IF NOT EXISTS short_code VARCHAR(8);

UPDATE storages
SET provider = COALESCE(provider, 'Railway');

UPDATE storages
SET path = COALESCE(path, object_key)
WHERE path IS NULL;

UPDATE storages
SET url = COALESCE(url, public_url)
WHERE url IS NULL;

UPDATE storages
SET file_size = COALESCE(file_size, size_bytes)
WHERE file_size IS NULL;

UPDATE storages
SET mime_type = COALESCE(mime_type, content_type)
WHERE mime_type IS NULL;

UPDATE storages
SET short_code = SUBSTRING(REPLACE(gen_random_uuid()::text, '-', '') FROM 1 FOR 8)
WHERE short_code IS NULL OR short_code = '';

ALTER TABLE storages
    ALTER COLUMN provider SET NOT NULL,
    ALTER COLUMN file_size SET NOT NULL,
    ALTER COLUMN short_code SET NOT NULL,
    ALTER COLUMN short_code TYPE VARCHAR(8);

DROP INDEX IF EXISTS storages_object_key_uq;
DROP INDEX IF EXISTS storages_deleted_at_idx;

ALTER TABLE storages
    DROP COLUMN IF EXISTS bucket_name,
    DROP COLUMN IF EXISTS object_key,
    DROP COLUMN IF EXISTS file_name,
    DROP COLUMN IF EXISTS content_type,
    DROP COLUMN IF EXISTS size_bytes,
    DROP COLUMN IF EXISTS size_text,
    DROP COLUMN IF EXISTS public_url,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS deleted_at;

CREATE UNIQUE INDEX IF NOT EXISTS storages_short_code_uq ON storages (short_code);
CREATE INDEX IF NOT EXISTS storages_provider_idx ON storages (provider);
CREATE INDEX IF NOT EXISTS storages_created_at_idx ON storages (created_at DESC);

COMMENT ON TABLE storages IS 'ตารางเก็บข้อมูลไฟล์ที่อัปโหลดขึ้น storage provider';
COMMENT ON COLUMN storages.id IS 'รหัสพื้นที่จัดเก็บ';
COMMENT ON COLUMN storages.provider IS 'ผู้ให้บริการ';
COMMENT ON COLUMN storages.path IS 'พาธที่จัดเก็บไฟล์';
COMMENT ON COLUMN storages.url IS 'ลิงก์เข้าถึงไฟล์';
COMMENT ON COLUMN storages.file_size IS 'ขนาดไฟล์ (Byte)';
COMMENT ON COLUMN storages.mime_type IS 'ประเภทไฟล์ (MIME Type)';
COMMENT ON COLUMN storages.created_at IS 'วันที่สร้างข้อมูล';
COMMENT ON COLUMN storages.short_code IS 'โค้ดสั้นสำหรับลิงก์สาธารณะของไฟล์ (6-8 ตัวอักษร)';
