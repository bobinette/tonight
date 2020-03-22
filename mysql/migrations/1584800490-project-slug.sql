-- Migration: project-slug
-- Created at: 2020-03-21 15:21:30
-- ====  UP  ====

BEGIN;

ALTER TABLE `projects`
    ADD COLUMN `slug` VARCHAR(256) NOT NULL  AFTER `name`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `projects`
    DROP COLUMN `slug`;

COMMIT;
