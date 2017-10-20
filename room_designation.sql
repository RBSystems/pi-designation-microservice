-- MySQL dump 10.16  Distrib 10.2.8-MariaDB, for osx10.12 (x86_64)
--
-- Host: localhost    Database: room_designation_2
-- ------------------------------------------------------
-- Server version	10.2.8-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `class_definitions`
--

DROP TABLE IF EXISTS `class_definitions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `class_definitions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` varchar(1024) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `class_definitions`
--

LOCK TABLES `class_definitions` WRITE;
/*!40000 ALTER TABLE `class_definitions` DISABLE KEYS */;
/*!40000 ALTER TABLE `class_definitions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `designation_definitions`
--

DROP TABLE IF EXISTS `designation_definitions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `designation_definitions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` varchar(1024) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `designation_definitions`
--

LOCK TABLES `designation_definitions` WRITE;
/*!40000 ALTER TABLE `designation_definitions` DISABLE KEYS */;
/*!40000 ALTER TABLE `designation_definitions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `microservice_definitions`
--

DROP TABLE IF EXISTS `microservice_definitions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `microservice_definitions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` varchar(1024) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `microservice_definitions`
--

LOCK TABLES `microservice_definitions` WRITE;
/*!40000 ALTER TABLE `microservice_definitions` DISABLE KEYS */;
INSERT INTO `microservice_definitions` VALUES (1,'configuration-database-microservice','fronts the av config db');
/*!40000 ALTER TABLE `microservice_definitions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `microservice_mappings`
--

DROP TABLE IF EXISTS `microservice_mappings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `microservice_mappings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `yaml` blob NOT NULL,
  `designation_id` int(11) NOT NULL,
  `class_id` int(11) NOT NULL,
  `microservice_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `designation_id` (`designation_id`,`class_id`,`microservice_id`,`yaml`(1024)),
  KEY `class_id` (`class_id`),
  KEY `microservice_id` (`microservice_id`),
  CONSTRAINT `microservice_mappings_ibfk_1` FOREIGN KEY (`designation_id`) REFERENCES `designation_definitions` (`id`),
  CONSTRAINT `microservice_mappings_ibfk_2` FOREIGN KEY (`class_id`) REFERENCES `class_definitions` (`id`),
  CONSTRAINT `microservice_mappings_ibfk_3` FOREIGN KEY (`microservice_id`) REFERENCES `microservice_definitions` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `microservice_mappings`
--

LOCK TABLES `microservice_mappings` WRITE;
/*!40000 ALTER TABLE `microservice_mappings` DISABLE KEYS */;
/*!40000 ALTER TABLE `microservice_mappings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rooms`
--

DROP TABLE IF EXISTS `rooms`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rooms` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `designation_id` int(11) NOT NULL,
  `ui_configuration` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `designation_id` (`designation_id`),
  CONSTRAINT `rooms_ibfk_1` FOREIGN KEY (`designation_id`) REFERENCES `designation_definitions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rooms`
--

LOCK TABLES `rooms` WRITE;
/*!40000 ALTER TABLE `rooms` DISABLE KEYS */;
/*!40000 ALTER TABLE `rooms` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `variable_definitions`
--

DROP TABLE IF EXISTS `variable_definitions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `variable_definitions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` varchar(1024) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `variable_definitions`
--

LOCK TABLES `variable_definitions` WRITE;
/*!40000 ALTER TABLE `variable_definitions` DISABLE KEYS */;
/*!40000 ALTER TABLE `variable_definitions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `variable_mappings`
--

DROP TABLE IF EXISTS `variable_mappings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `variable_mappings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `value` varchar(512) NOT NULL,
  `designation_id` int(11) NOT NULL,
  `class_id` int(11) NOT NULL,
  `variable_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `designation_id` (`designation_id`,`class_id`,`variable_id`,`value`),
  KEY `class_id` (`class_id`),
  KEY `variable_id` (`variable_id`),
  CONSTRAINT `variable_mappings_ibfk_1` FOREIGN KEY (`designation_id`) REFERENCES `designation_definitions` (`id`),
  CONSTRAINT `variable_mappings_ibfk_2` FOREIGN KEY (`class_id`) REFERENCES `class_definitions` (`id`),
  CONSTRAINT `variable_mappings_ibfk_3` FOREIGN KEY (`variable_id`) REFERENCES `variable_definitions` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `variable_mappings`
--

LOCK TABLES `variable_mappings` WRITE;
/*!40000 ALTER TABLE `variable_mappings` DISABLE KEYS */;
/*!40000 ALTER TABLE `variable_mappings` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-10-19 21:03:06
