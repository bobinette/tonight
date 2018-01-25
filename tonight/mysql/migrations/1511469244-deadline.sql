-- Migration: deadline
-- Created at: 2017-11-23 21:34:04
-- ====  UP  ====

ALTER TABLE `tasks`
    ADD COLUMN `deadline` DATE DEFAULT NULL AFTER `duration`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `deadline`;

COMMIT;
