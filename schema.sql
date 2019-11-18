DROP DATABASE IF EXISTS rss_reader;
CREATE DATABASE rss_reader;
USE rss_reader;

CREATE TABLE users (
  id                INT AUTO_INCREMENT,
  email             VARCHAR(100) NOT NULL,
  password_salt     VARCHAR(250) NOT NULL,
  password_hash     VARCHAR(250) NOT NULL,
  create_time       TIME NOT NULL,
  update_time       TIME NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT unique_email UNIQUE (email)
);

CREATE TABLE sources (
  id                INT AUTO_INCREMENT,
  user_id           INT NOT NULL,
  link              VARCHAR(250) NOT NULL,
  create_time       TIME NOT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE items (
  id              INT AUTO_INCREMENT,
  source_id       INT NOT NULL,
  title           VARCHAR(250) NOT NULL,
  link            VARCHAR(250) NOT NULL,
  description     VARCHAR(400),
  pubDate         TIME NOT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY(source_id) REFERENCES sources(id)
);

