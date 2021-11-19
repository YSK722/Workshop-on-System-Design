INSERT INTO `tasks` (`title`) VALUES ("sample-task-01");
INSERT INTO `tasks` (`title`) VALUES ("sample-task-02");
INSERT INTO `tasks` (`title`) VALUES ("sample-task-03");
INSERT INTO `tasks` (`title`) VALUES ("sample-task-04");
INSERT INTO `tasks` (`title`, `is_done`) VALUES ("sample-task-05", true);
UPDATE `tasks` SET `deadline`=DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 1 YEAR);