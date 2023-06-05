CREATE TABLE `test`.`users` (
  `id` INT NULL AUTO_INCREMENT,
  `email` VARCHAR(50) NOT NULL,
  `password` MEDIUMTEXT NOT NULL,
  `created` TIMESTAMP(2) NOT NULL,
  `last_login` TIMESTAMP(2) NOT NULL,
  PRIMARY KEY (`email`),
  UNIQUE INDEX `userscol_UNIQUE` (`userscol` ASC) VISIBLE);
