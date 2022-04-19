CREATE TABLE socket_types (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

INSERT INTO socket_types (name)
    VALUES ("unknown"), ('air conditioner'), ('heater'), ('generator');

CREATE TABLE sockets (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    mib_address VARCHAR(100) NOT NULL,
	netping_address VARCHAR(50) NOT NULL,
    socket_type_id INT NOT NULL,
    FOREIGN KEY (socket_type_id) REFERENCES socket_types(id)
);
