-- MySQL dump 10.16  Distrib 10.2.8-MariaDB, for osx10.12 (x86_64)
--
-- Host: localhost    Database: room_designation
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
-- Table structure for table `designation_definition`
--

DROP TABLE IF EXISTS `designation_definition`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `designation_definition` (
  `designation` varchar(20) NOT NULL,
  `designation_ID` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`designation_ID`),
  UNIQUE KEY `designation` (`designation`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `designation_definition`
--

LOCK TABLES `designation_definition` WRITE;
/*!40000 ALTER TABLE `designation_definition` DISABLE KEYS */;
INSERT INTO `designation_definition` VALUES ('development',1),('production',4),('stage',3),('testing',2);
/*!40000 ALTER TABLE `designation_definition` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rooms`
--

DROP TABLE IF EXISTS `rooms`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rooms` (
  `room_ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL,
  `designation_ID` int(11) DEFAULT NULL,
  `ui_config` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  PRIMARY KEY (`room_ID`),
  KEY `designation_ID` (`designation_ID`),
  CONSTRAINT `rooms_ibfk_1` FOREIGN KEY (`designation_ID`) REFERENCES `designation_definition` (`designation_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rooms`
--

LOCK TABLES `rooms` WRITE;
/*!40000 ALTER TABLE `rooms` DISABLE KEYS */;
INSERT INTO `rooms` VALUES (1,'ITB-1101',1,'{\'c\':\'Monster\'}'),(2,'ITB-1102',1,'{\'d\':\'Monster\'}'),(3,'ITB-3006',1,'{\"apiconfig\":{\"enabled\":true,\"backups\":{\"1\":\"ITB-3006-CP2\"}},\"devices\":[{\"ui\":\"circle-default\",\"inputdevices\":[{\"name\":\"appleTV\",\"icon\":\"apple\"}],\"displays\":[],\"audio\":[],\"features\":[]}]}'),(4,'CTB-410',1,'{\"apiconfig\":{\"enabled\":true,\"backups\":{\"1\":\"ITB-3006-CP2\"}},\"devices\":[{\"ui\":\"circle-default\",\"inputdevices\":[{\"name\":\"appleTV\",\"icon\":\"apple\"}],\"displays\":[],\"audio\":[],\"features\":[]}]}');
/*!40000 ALTER TABLE `rooms` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ui_config`
--

DROP TABLE IF EXISTS `ui_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ui_config` (
  `config_ID` int(11) NOT NULL AUTO_INCREMENT,
  `room_ID` int(11) NOT NULL,
  `config_file` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  PRIMARY KEY (`config_ID`),
  KEY `room_ID` (`room_ID`),
  CONSTRAINT `ui_config_ibfk_1` FOREIGN KEY (`room_ID`) REFERENCES `rooms` (`room_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ui_config`
--

LOCK TABLES `ui_config` WRITE;
/*!40000 ALTER TABLE `ui_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `ui_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `variables`
--

DROP TABLE IF EXISTS `variables`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `variables` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `desig_id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `value` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `desig_name` (`name`,`desig_id`),
  KEY `desig_id` (`desig_id`),
  CONSTRAINT `variables_ibfk_1` FOREIGN KEY (`desig_id`) REFERENCES `designation_definition` (`designation_ID`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `variables`
--

LOCK TABLES `variables` WRITE;
/*!40000 ALTER TABLE `variables` DISABLE KEYS */;
INSERT INTO `variables` VALUES (5,2,'CONFIGURATION_DATABASE_PASSWORD','eMonster'),(7,3,'CONFIGURATION_DATABASE_PASSWORD','eMonster');
/*!40000 ALTER TABLE `variables` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-10-03 16:15:23
