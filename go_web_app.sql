CREATE TABLE user (
    id CHAR(36) NOT NULL PRIMARY KEY,
    username VARCHAR(20) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password BINARY(64) NOT NULL,
    salt BINARY(16) NOT NULL,
    registration_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_account_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    is_account_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE account_confirmation (
    user_id CHAR(36) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    confirmation_code CHAR(36) UNIQUE NOT NULL,
    security_code CHAR(6) NOT NULL
);

CREATE TABLE chat (
	id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(100)
);

CREATE TABLE chat_participant (
	chat_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
	FOREIGN KEY (chat_id) REFERENCES chat(id),
    FOREIGN KEY (user_id) REFERENCES user(id)
);

CREATE TABLE message (
	id CHAR(36) NOT NULL PRIMARY KEY,
    chat_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chat(id),
    FOREIGN KEY (user_id) REFERENCES user(id),
    content VARCHAR(500) NOT NULL,
    created DATETIME NOT NULL
);



drop table user;
drop table account_confirmation;
select * from user;