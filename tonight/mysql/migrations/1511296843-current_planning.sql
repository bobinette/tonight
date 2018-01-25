-- Migration: current_planning
-- Created at: 2017-11-21 21:40:43
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `planning` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,

    `duration` VARCHAR(64) NOT NULL, -- duration as a string
    `dismissed` BOOLEAN NOT NULL,
    `startedAt` DATETIME NOT NULL,

    PRIMARY KEY(`id`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS `planning_has_task` (
    `planning_id` BIGINT NOT NULL,
    `rank` INT NOT NULL,
    `task_id` BIGINT NOT NULL,

    PRIMARY KEY(`planning_id`, `rank`),

    CONSTRAINT `fk_planning_has_task_planning` FOREIGN KEY (`planning_id`) REFERENCES `planning` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_planning_has_task_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

COMMIT;

-- ==== DOWN ====

BEGIN;

DROP TABLE IF EXISTS `planning_has_task`;
DROP TABLE IF EXISTS `planning`;

COMMIT;
