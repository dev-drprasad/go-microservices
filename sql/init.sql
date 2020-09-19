-- CREATE DOMAIN uint AS int CHECK(VALUE >= 0);
CREATE EXTENSION IF NOT EXISTS pgcrypto;


CREATE TYPE userRole AS ENUM('staff', 'admin');
CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  name varchar(60) NOT NULL,
  username varchar(20) UNIQUE NOT NULL,
  password varchar(72) NOT NULL,
  role userRole NOT NULL DEFAULT 'staff'
);
