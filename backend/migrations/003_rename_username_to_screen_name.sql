-- Migration: rename username → screen_name (aligns DB with API/domain model)
-- HB-AUTH-02

ALTER TABLE IF EXISTS users RENAME COLUMN username TO screen_name;

-- Update full-text search index to use new column name
-- (the tsvector expression in queries will reference screen_name instead)
