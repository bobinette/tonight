-- Migration: planning_search_query
-- Created at: 2017-12-23 12:09:11
-- ====  UP  ====

BEGIN;

ALTER TABLE `planning`
    ADD COLUMN `search_query` TEXT NOT NULL AFTER `strict`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `planning`
    DROP COLUMN `search_query`;

COMMIT;
