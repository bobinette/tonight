-- Migration: users
-- Created at: 2017-12-01 21:02:10
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `user` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(256) NOT NULL,

    `created_at` DATETIME NOT NULL,

    PRIMARY KEY(`id`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

COMMIT;

-- ==== DOWN ====

BEGIN;

DROP TABLE IF EXISTS `user`;

COMMIT;
