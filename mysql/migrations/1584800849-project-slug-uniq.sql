-- Migration: project-slug-uniq
-- Created at: 2020-03-21 15:27:29
-- ====  UP  ====

BEGIN;

ALTER TABLE `projects`
    ADD CONSTRAINT `u_project_slug` UNIQUE KEY(`slug`);

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `projects`
    DROP INDEX `u_project_slug`;


COMMIT;
