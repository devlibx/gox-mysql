#DROP TABLE `users`;
CREATE TABLE IF NOT EXISTS `integrating_tests_users`
(
    `id`      int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name`    varchar(265) NOT NULL DEFAULT '',
    `deleted` int          DEFAULT 0,
    PRIMARY KEY (`id`)
);