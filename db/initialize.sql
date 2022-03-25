CREATE USER `shion`@`%` IDENTIFIED BY 'password';
GRANT INSERT,SELECT,UPDATE,DELETE ON `shionstagram_db`.* TO `shion`@`%`;

CREATE DATABASE IF NOT EXISTS `shionstagram_db`;

CREATE TABLE IF NOT EXISTS `shionstagram_db`.`posts` (
    `id`       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `uuid`     VARBINARY(16) NOT NULL,
    `name`     VARCHAR(255) NOT NULL,
    `location` VARCHAR(255) NOT NULL,
    `message`  VARCHAR(2048) NOT NULL,
    `img_src`  VARCHAR(255),
    `avatar`   TINYINT UNSIGNED NOT NULL,
    `pending`  BOOLEAN NOT NULL,

    PRIMARY KEY ( `id` )
)
