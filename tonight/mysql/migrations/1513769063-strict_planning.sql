-- Migration: strict_planning
-- Created at: 2017-12-20 12:24:23
-- ====  UP  ====

BEGIN;

ALTER TABLE `planning`
    ADD COLUMN `strict` BOOLEAN NOT NULL AFTER `duration`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `planning`
    DROP COLUMN `strict`;

COMMIT;
