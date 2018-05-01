-- Migration: remove-postponed-column
-- Created at: 2018-05-01 14:43:37
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `postponed_until`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `postponed_until` DATE DEFAULT NULL AFTER `deadline`;

COMMIT;
