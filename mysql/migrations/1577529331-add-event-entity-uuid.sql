-- Migration: add-event-entity-uuid
-- Created at: 2019-12-28 11:35:31
-- ====  UP  ====

BEGIN;

ALTER TABLE `events`
    ADD COLUMN `entity_uuid` VARCHAR(36) NOT NULL AFTER `type`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `events`
    DROP COLUMN `entity_uuid`;

COMMIT;
