-- Migration: 006_enable_pg_trgm
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Description: Enable pg_trgm extension for trigram text similarity and fuzzy search

