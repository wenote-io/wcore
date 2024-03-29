CREATE DATABASE wcore
USE wcore
CREATE TABLE w_user(
    id INT NOT NULL PRIMARY KEY auto_increment,  -- primary key column
    w_id  VARCHAR(32) NOT NULL unique,
    ext_id VARCHAR(32) NOT NULL,
    w_type       INT NOT NULL DEFAULT 0,
    create_time BIGINT NOT NULL DEFAULT 0,
    delete_time BIGINT NOT NULL DEFAULT 0,
    update_time BIGINT NOT NULL DEFAULT 0
);  

CREATE TABLE user_action_stat(
    id INT NOT NULL PRIMARY KEY auto_increment, -- primary key column
    w_id VARCHAR(32) NOT NULL ,
    act_type INT NOT NULL DEFAULT 0,
    act_val  BIGINT NOT NULL DEFAULT 0,
    act_unit VARCHAR(32) NOT NULL DEFAULT "no",
    create_time BIGINT NOT NULL DEFAULT 0,
    delete_time BIGINT NOT NULL DEFAULT 0,
    update_time BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE w_note(
    id INT NOT NULL PRIMARY KEY auto_increment, -- primary key column
    w_id VARCHAR(32) NOT NULL ,
    n_id VARCHAR(32) NOT NULL ,
    w_mood INT NOT NULL,
    w_desc TEXT NOT NULL,
    w_action_type INT NOT NULL DEFAULT 0, 
    create_time BIGINT NOT NULL DEFAULT 0,
    delete_time BIGINT NOT NULL DEFAULT 0,
    update_time BIGINT NOT NULL DEFAULT 0
);  