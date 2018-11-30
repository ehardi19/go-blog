CREATE DATABASE IF NOT EXISTS `mydb`
USE `mydb`;

CREATE TABLE IF NOT EXISTS  `users` (
  username text PRIMARY KEY,
  password text
);

CREATE TABLE IF NOT EXIST `article` (
  id text PRIMARY KEY,
  title text,
  content text
)