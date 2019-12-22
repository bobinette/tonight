-- Migration: event-table
-- Created at: 2019-11-23 10:47:43
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `events` (
    `uuid` VARCHAR(36) NOT NULL,

    `type` VARCHAR(256) NOT NULL,
    `payload` MEDIUMTEXT NOT NULL,

    `created_at` DATETIME NOT NULL,

    PRIMARY KEY (`uuid`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

COMMIT;

-- ==== DOWN ====

BEGIN;

DROP TABLE IF EXISTS `events`;

COMMIT;
