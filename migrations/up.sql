-- These migrations could not be automated
-- since this application will use already existing db
-- this file is here just to make stuff more clear

CREATE TABLE `netping_list` (
  `id` int NOT NULL,
  `firmware_id` int NOT NULL,
  `name` varchar(50) NOT NULL,
  `domain` varchar(50) NOT NULL,
  `host` varchar(100) DEFAULT NULL,
  `web` varchar(255) DEFAULT NULL,
  `position` int NOT NULL,
  `diesel` tinyint NOT NULL,
  `power_note` varchar(50) DEFAULT NULL,
  `sensors` text,
  `devices` text,
  `request` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

CREATE TABLE socket_types (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

INSERT INTO socket_types (name)
    VALUES ("unknown"), ('air conditioner'), ('heater'), ('generator');

CREATE TABLE sockets (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    mib_address VARCHAR(100) NOT NULL,
	netping_id INT NOT NULL,
    socket_type_id INT NOT NULL,
    FOREIGN KEY (socket_type_id) REFERENCES socket_types(id),
    FOREIGN KEY (netping_id) REFERENCES netping_list(id)
);

CREATE TABLE netping_manager_users (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    is_admin BOOL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

INSERT INTO netping_manager_users(name, password, is_admin)
    VALUES('admin', 'button_lord', true);