-- Migration: task-table
-- Created at: 2019-11-23 10:40:10
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `tasks` (
    `uuid` VARCHAR(36) NOT NULL,

    `title` TEXT NOT NULL,

    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,

    PRIMARY KEY (`uuid`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

COMMIT;

-- ==== DOWN ====

BEGIN;

DROP TABLE IF EXISTS `tasks`;

COMMIT;
