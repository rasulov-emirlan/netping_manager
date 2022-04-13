-- phpMyAdmin SQL Dump
-- version 5.1.0
-- https://www.phpmyadmin.net/
--
-- Хост: localhost
-- Время создания: Апр 13 2022 г., 13:14
-- Версия сервера: 5.7.33-0ubuntu0.16.04.1
-- Версия PHP: 7.4.16

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- База данных: `npLog`
--
CREATE DATABASE IF NOT EXISTS `npLog` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `npLog`;

-- --------------------------------------------------------

--
-- Структура таблицы `netping_history`
--

CREATE TABLE `netping_history` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `netping_id` int(11) NOT NULL,
  `snapshot` text NOT NULL,
  `create_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

--
-- Индексы сохранённых таблиц
--

--
-- Индексы таблицы `netping_history`
--
ALTER TABLE `netping_history`
  ADD PRIMARY KEY (`id`),
  ADD KEY `netping_id` (`netping_id`);

--
-- AUTO_INCREMENT для сохранённых таблиц
--

--
-- AUTO_INCREMENT для таблицы `netping_history`
--
ALTER TABLE `netping_history`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- База данных: `nrp`
--
CREATE DATABASE IF NOT EXISTS `nrp` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `nrp`;

-- --------------------------------------------------------

--
-- Структура таблицы `inode_list`
--

CREATE TABLE `inode_list` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `domain` varchar(100) NOT NULL,
  `host` varchar(15) NOT NULL,
  `model` varchar(100) NOT NULL,
  `desc` varchar(255) NOT NULL,
  `request` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `netping_diesel`
--

CREATE TABLE `netping_diesel` (
  `id` bigint(20) NOT NULL,
  `netping_id` int(11) NOT NULL,
  `ontime` timestamp NULL DEFAULT NULL,
  `offtime` timestamp NULL DEFAULT NULL,
  `worktime` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `netping_diesel_tmp`
--

CREATE TABLE `netping_diesel_tmp` (
  `netping_id` int(11) NOT NULL,
  `ontime` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `netping_events`
--

CREATE TABLE `netping_events` (
  `uid` varchar(50) NOT NULL,
  `netping_id` int(11) NOT NULL,
  `source` varchar(50) NOT NULL,
  `code` varchar(50) NOT NULL,
  `value` int(11) DEFAULT NULL,
  `create_date` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `netping_firmware`
--

CREATE TABLE `netping_firmware` (
  `id` int(11) NOT NULL,
  `fw_name` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `netping_list`
--

CREATE TABLE `netping_list` (
  `id` int(11) NOT NULL,
  `firmware_id` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `domain` varchar(50) NOT NULL,
  `host` varchar(100) DEFAULT NULL,
  `web` varchar(255) DEFAULT NULL,
  `position` int(11) NOT NULL,
  `diesel` tinyint(4) NOT NULL,
  `power_note` varchar(50) DEFAULT NULL,
  `sensors` text,
  `devices` text,
  `request` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_alarms`
--

CREATE TABLE `pl_alarms` (
  `id` int(11) NOT NULL,
  `code` int(11) NOT NULL,
  `char_value` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `note` text
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_directions`
--

CREATE TABLE `pl_directions` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `code` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_interfaces`
--

CREATE TABLE `pl_interfaces` (
  `id` int(11) NOT NULL,
  `number` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_lines`
--

CREATE TABLE `pl_lines` (
  `id` int(11) NOT NULL,
  `selection_id` int(11) NOT NULL,
  `pl_id_a` int(11) NOT NULL,
  `pl_id_b` int(11) NOT NULL,
  `port_a` int(11) NOT NULL,
  `port_b` int(11) NOT NULL,
  `lambda` int(11) NOT NULL,
  `attn` varchar(50) DEFAULT NULL,
  `tx` varchar(50) DEFAULT NULL,
  `rx` varchar(50) DEFAULT NULL,
  `fecr` varchar(50) DEFAULT NULL,
  `osnr` varchar(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_links`
--

CREATE TABLE `pl_links` (
  `id` int(11) NOT NULL,
  `from` varchar(100) NOT NULL,
  `to` varchar(100) NOT NULL,
  `number` int(11) NOT NULL,
  `direction` enum('fwd','bwd') NOT NULL,
  `length` float NOT NULL,
  `signal` float NOT NULL,
  `norm` float NOT NULL,
  `desc` varchar(100) DEFAULT NULL,
  `port` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_modules`
--

CREATE TABLE `pl_modules` (
  `id` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `number` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_packetlights`
--

CREATE TABLE `pl_packetlights` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(50) NOT NULL,
  `domain` varchar(100) NOT NULL,
  `ip` varchar(50) NOT NULL,
  `old_ip` varchar(50) NOT NULL,
  `direction` varchar(255) NOT NULL,
  `modules` text,
  `links` text,
  `syslog` tinyint(1) NOT NULL DEFAULT '1',
  `request` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_ports`
--

CREATE TABLE `pl_ports` (
  `id` int(11) NOT NULL,
  `number` int(11) NOT NULL,
  `tn_name` varchar(50) DEFAULT NULL,
  `il_name` varchar(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_port_channel`
--

CREATE TABLE `pl_port_channel` (
  `id` int(11) NOT NULL,
  `number` int(11) NOT NULL,
  `sfp_index` int(11) NOT NULL,
  `fec_index` int(11) NOT NULL,
  `osn_index` int(11) DEFAULT '0',
  `desc` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_preamp`
--

CREATE TABLE `pl_preamp` (
  `id` int(11) NOT NULL,
  `domain` varchar(100) NOT NULL,
  `number` tinyint(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_selections`
--

CREATE TABLE `pl_selections` (
  `id` int(11) NOT NULL,
  `direction_id` int(11) NOT NULL,
  `route` tinyint(4) NOT NULL,
  `name` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `pl_types`
--

CREATE TABLE `pl_types` (
  `id` int(11) NOT NULL,
  `code` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `roles`
--

CREATE TABLE `roles` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `desc` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Индексы сохранённых таблиц
--

--
-- Индексы таблицы `inode_list`
--
ALTER TABLE `inode_list`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `netping_diesel`
--
ALTER TABLE `netping_diesel`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `netping_diesel_tmp`
--
ALTER TABLE `netping_diesel_tmp`
  ADD PRIMARY KEY (`netping_id`);

--
-- Индексы таблицы `netping_events`
--
ALTER TABLE `netping_events`
  ADD PRIMARY KEY (`uid`),
  ADD KEY `netping_id` (`netping_id`);

--
-- Индексы таблицы `netping_firmware`
--
ALTER TABLE `netping_firmware`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `netping_list`
--
ALTER TABLE `netping_list`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `pl_alarms`
--
ALTER TABLE `pl_alarms`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `pl_directions`
--
ALTER TABLE `pl_directions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `code` (`code`);

--
-- Индексы таблицы `pl_interfaces`
--
ALTER TABLE `pl_interfaces`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `pl_lines`
--
ALTER TABLE `pl_lines`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `pl_links`
--
ALTER TABLE `pl_links`
  ADD PRIMARY KEY (`id`),
  ADD KEY `from` (`from`),
  ADD KEY `to` (`to`);

--
-- Индексы таблицы `pl_modules`
--
ALTER TABLE `pl_modules`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `pl_packetlights`
--
ALTER TABLE `pl_packetlights`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `domain` (`domain`),
  ADD KEY `fk_type_code` (`type`),
  ADD KEY `fk_direction_code` (`direction`);

--
-- Индексы таблицы `pl_ports`
--
ALTER TABLE `pl_ports`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `pl_port_channel`
--
ALTER TABLE `pl_port_channel`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `pl_preamp`
--
ALTER TABLE `pl_preamp`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `unique_index` (`domain`,`number`);

--
-- Индексы таблицы `pl_selections`
--
ALTER TABLE `pl_selections`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `unique_index` (`name`,`direction_id`,`route`),
  ADD KEY `direction` (`direction_id`);

--
-- Индексы таблицы `pl_types`
--
ALTER TABLE `pl_types`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `code` (`code`);

--
-- Индексы таблицы `roles`
--
ALTER TABLE `roles`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT для сохранённых таблиц
--

--
-- AUTO_INCREMENT для таблицы `inode_list`
--
ALTER TABLE `inode_list`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `netping_diesel`
--
ALTER TABLE `netping_diesel`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `netping_firmware`
--
ALTER TABLE `netping_firmware`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `netping_list`
--
ALTER TABLE `netping_list`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_alarms`
--
ALTER TABLE `pl_alarms`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_directions`
--
ALTER TABLE `pl_directions`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_interfaces`
--
ALTER TABLE `pl_interfaces`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_lines`
--
ALTER TABLE `pl_lines`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_links`
--
ALTER TABLE `pl_links`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_modules`
--
ALTER TABLE `pl_modules`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_packetlights`
--
ALTER TABLE `pl_packetlights`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_ports`
--
ALTER TABLE `pl_ports`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_port_channel`
--
ALTER TABLE `pl_port_channel`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_preamp`
--
ALTER TABLE `pl_preamp`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_selections`
--
ALTER TABLE `pl_selections`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `pl_types`
--
ALTER TABLE `pl_types`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `roles`
--
ALTER TABLE `roles`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- Ограничения внешнего ключа сохраненных таблиц
--

--
-- Ограничения внешнего ключа таблицы `netping_events`
--
ALTER TABLE `netping_events`
  ADD CONSTRAINT `netping_events_ibfk_1` FOREIGN KEY (`netping_id`) REFERENCES `netping_list` (`id`);

--
-- Ограничения внешнего ключа таблицы `pl_links`
--
ALTER TABLE `pl_links`
  ADD CONSTRAINT `pl_links_ibfk_1` FOREIGN KEY (`from`) REFERENCES `pl_packetlights` (`domain`) ON UPDATE CASCADE,
  ADD CONSTRAINT `pl_links_ibfk_2` FOREIGN KEY (`to`) REFERENCES `pl_packetlights` (`domain`) ON UPDATE CASCADE;

--
-- Ограничения внешнего ключа таблицы `pl_packetlights`
--
ALTER TABLE `pl_packetlights`
  ADD CONSTRAINT `fk_direction_code` FOREIGN KEY (`direction`) REFERENCES `pl_directions` (`code`) ON UPDATE CASCADE,
  ADD CONSTRAINT `fk_type_code` FOREIGN KEY (`type`) REFERENCES `pl_types` (`code`);

--
-- Ограничения внешнего ключа таблицы `pl_preamp`
--
ALTER TABLE `pl_preamp`
  ADD CONSTRAINT `pl_preamp_ibfk_1` FOREIGN KEY (`domain`) REFERENCES `pl_packetlights` (`domain`);
--
-- База данных: `plSyslog`
--
CREATE DATABASE IF NOT EXISTS `plSyslog` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `plSyslog`;

-- --------------------------------------------------------

--
-- Структура таблицы `plog`
--

CREATE TABLE `plog` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `fromhostip` varchar(50) DEFAULT NULL,
  `rawmsg` text,
  `msg` text,
  `timereceipt` varchar(50) DEFAULT NULL,
  `timereported` datetime DEFAULT NULL,
  `timegenerated` datetime DEFAULT NULL,
  `timerecord` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

--
-- Индексы сохранённых таблиц
--

--
-- Индексы таблицы `plog`
--
ALTER TABLE `plog`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fromhostip` (`fromhostip`),
  ADD KEY `timereceipt` (`timereceipt`);
ALTER TABLE `plog` ADD FULLTEXT KEY `rawmsg` (`rawmsg`);

--
-- AUTO_INCREMENT для сохранённых таблиц
--

--
-- AUTO_INCREMENT для таблицы `plog`
--
ALTER TABLE `plog`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
