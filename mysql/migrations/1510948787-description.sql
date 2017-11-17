-- Migration: description
-- Created at: 2017-11-17 20:59:47
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `description` TEXT NOT NULL AFTER `title`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `description`;

COMMIT;
