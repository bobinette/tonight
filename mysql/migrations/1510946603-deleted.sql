-- Migration: deleted
-- Created at: 2017-11-17 20:23:23
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `deleted` BOOLEAN DEFAULT FALSE AFTER `updated_at`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `deleted`;

COMMIT;
