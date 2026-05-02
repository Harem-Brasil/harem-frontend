-- URLs de mídia no catálogo + índice por criador e visibilidade (listagens filtradas).

ALTER TABLE creator_catalog
    ADD COLUMN IF NOT EXISTS media_urls TEXT[] NOT NULL DEFAULT '{}';

CREATE INDEX IF NOT EXISTS idx_creator_catalog_creator_visibility_active
    ON creator_catalog (creator_id, visibility)
    WHERE deleted_at IS NULL;
