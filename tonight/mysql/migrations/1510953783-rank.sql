-- Migration: rank
-- Created at: 2017-11-17 22:23:03
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `rank` INTEGER NOT NULL AFTER `description`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `rank`;

COMMIT;

