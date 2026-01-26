-- Migration: 006_enable_pg_trgm
-- Description: Rollback pg_trgm extension

DROP EXTENSION IF EXISTS pg_trgm;

