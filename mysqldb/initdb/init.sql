-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
-- -----------------------------------------------------
-- Schema mysql_ztb
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema mysql_ztb
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `mysql_ztb` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci ;
USE `mysql_ztb` ;

-- -----------------------------------------------------
-- Table `mysql_ztb`.`USER`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mysql_ztb`.`USER` (
    `id` VARCHAR(100) NOT NULL,
    `first_name` VARCHAR(25) NOT NULL,
    `last_name` VARCHAR(40) NOT NULL,
    `email` VARCHAR(50) NOT NULL,
    `username` VARCHAR(45) NOT NULL,
    `is_admin` TINYINT(1) NOT NULL,
    `is_enabled` TINYINT(1) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `email` (`email` ASC) VISIBLE,
    UNIQUE INDEX `username` (`username` ASC) VISIBLE)
    ENGINE = InnoDB
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `mysql_ztb`.`NOTE`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mysql_ztb`.`NOTE` (
    `id` VARCHAR(100) NOT NULL,
    `title` VARCHAR(200) NOT NULL,
    `content` LONGTEXT NOT NULL,
    `date_of_creation` DATETIME NOT NULL,
    `date_of_modification` DATETIME NULL DEFAULT NULL,
    `is_shared` TINYINT(1) NOT NULL,
    `id_user` VARCHAR(100) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `title` (`title` ASC) VISIBLE,
    INDEX `id_user` (`id_user` ASC) VISIBLE,
    CONSTRAINT `NOTE_ibfk_1`
    FOREIGN KEY (`id_user`)
    REFERENCES `mysql_ztb`.`USER` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
    ENGINE = InnoDB
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_0900_ai_ci;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
