-- Migration: project-description
-- Created at: 2020-03-22 18:45:10
-- ====  UP  ====

BEGIN;

ALTER TABLE `projects`
    ADD COLUMN `description` TEXT NOT NULL  AFTER `name`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `projects`
    DROP COLUMN `description`;

COMMIT;
