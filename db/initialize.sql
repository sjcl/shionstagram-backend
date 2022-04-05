CREATE USER `shion`@`%` IDENTIFIED BY 'password';
GRANT INSERT,SELECT,UPDATE,DELETE ON `shionstagram_db`.* TO `shion`@`%`;

CREATE DATABASE IF NOT EXISTS `shionstagram_db`;

CREATE TABLE IF NOT EXISTS `shionstagram_db`.`messages` (
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `uuid`         UUID    NOT NULL,
    `twitter_name` VARCHAR(15)      NOT NULL,
    `name`         VARCHAR(100)     NOT NULL,
    `location`     VARCHAR(100),
    `message`      VARCHAR(500)     NOT NULL,
    `image`        VARCHAR(255),
    `avatar`       TINYINT UNSIGNED NOT NULL,
    `is_pending`   BOOLEAN          NOT NULL DEFAULT 1,
    `created_at`   TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ( `id` )
)
