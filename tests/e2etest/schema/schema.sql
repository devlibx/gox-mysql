DROP TABLE IF EXISTS `integrating_tests_users`;
CREATE TABLE IF NOT EXISTS `integrating_tests_users`
(
    `id`      int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name`    varchar(265) NOT NULL DEFAULT '',
    `department` varchar(265) NOT NULL DEFAULT 'na',
    `deleted` int          DEFAULT 0,
    PRIMARY KEY (`id`)
);