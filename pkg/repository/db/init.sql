CREATE TABLE users (
    id SERIAL,
    login VARCHAR(36),
    password VARCHAR(36)
);

CREATE TABLE adverts (
    id SERIAL,
	user_id INT,
    header VARCHAR(64),
    text VARCHAR(512),
    address VARCHAR(64),
    image_url VARCHAR(64),
    price REAL
);


