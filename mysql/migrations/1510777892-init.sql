-- Migration: init
-- Created at: 2017-11-15 21:31:32
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `tasks` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,

    `title` VARCHAR(2048) NOT NULL,
    `done` BOOLEAN NOT NULL,

    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,

    PRIMARY KEY (`id`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

COMMIT;

-- ==== DOWN ====

BEGIN;

COMMIT;
