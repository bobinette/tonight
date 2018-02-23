-- Migration: add-postponed-until
-- Created at: 2018-02-23 18:34:06
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `postponed_until` DATE DEFAULT NULL AFTER `deadline`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `postponed_until`;

COMMIT;
