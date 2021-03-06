
-- create database ustapi;
-- CREATE USER user WITH PASSWORD 'password';
-- GRANT ALL PRIVILEGES ON DATABASE ustapi TO ustapi;

CREATE DOMAIN uint AS int CHECK(VALUE >= 0);
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION triggerSetUpdatedAt()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updatedAt = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS organizations (
  id serial PRIMARY KEY,
  createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name varchar(60) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS branches (
  id serial PRIMARY KEY,
  createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name varchar(60) NOT NULL,
  organizationId int NOT NULL,
  address varchar(255) NOT NULL,
  phoneNumber varchar(12) NOT NULL,
  FOREIGN KEY (organizationId) REFERENCES organizations (id)
);
-- ALTER TABLE branches ADD COLUMN address varchar(255) NOT NULL DEFAULT '';
-- ALTER TABLE branches ADD COLUMN phoneNumber varchar(12) NOT NULL DEFAULT '0000000000';



CREATE TYPE userRole AS ENUM('staff', 'admin', 'superadmin');
-- ALTER TYPE userRole ADD VALUE 'superadmin';
CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name varchar(30) NOT NULL,
  username varchar(20) UNIQUE NOT NULL,
  password varchar(72) NOT NULL,
  role userRole NOT NULL DEFAULT 'staff'
);
ALTER TABLE users ADD COLUMN branchId int NOT NULL REFERENCES branches(id);
-- ALTER TABLE users ADD COLUMN organizationId int NOT NULL REFERENCES organizations(id);

INSERT INTO organizations (name) VALUES ('Super Admin');
INSERT INTO branches (name, organizationId, address, phoneNumber) VALUES ('Super Admin', 1, '', '');
INSERT INTO users (name, username, password, branchId, role) VALUES ('Super Admin', 'superadmin', crypt('00000', gen_salt('bf')), 1, 'superadmin');

CREATE TABLE IF NOT EXISTS brands (
  id serial PRIMARY KEY,
  createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name varchar(20) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
  id serial PRIMARY KEY,
  createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name varchar(60) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
  id serial PRIMARY KEY,
  createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name varchar(60) NOT NULL,
  cost numeric(8,2) NOT NULL,
  sellPrice numeric(8,2) NOT NULL,
  brandId int NOT NULL,
  categoryId int NOT NULL,
  imageUrls varchar(255) array[10] NOT NULL,
  stock uint NOT NULL,
  FOREIGN KEY (brandId) REFERENCES brands(id),
  FOREIGN KEY (categoryId) REFERENCES categories(id)
);

CREATE TABLE IF NOT EXISTS product_images (
  id serial PRIMARY KEY,
  createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  imageUrl varchar(255) UNIQUE NOT NULL,
  productId uint NULL,
  FOREIGN KEY (productId) REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS customers (
  id SERIAL NOT NULL PRIMARY KEY,
  createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name varchar(30) NOT NULL,
  address varchar(255) NOT NULL,
  zipcode varchar(10) NOT NULL,
  phoneNumber varchar(15) NOT NULL
);

CREATE TRIGGER setUpdatedAt
BEFORE UPDATE ON customers
FOR EACH ROW
EXECUTE PROCEDURE triggerSetUpdatedAt();

CREATE TYPE orderStatus AS ENUM('new', 'preparation', 'ready', 'delivered');

CREATE TABLE IF NOT EXISTS orders (
  id SERIAL NOT NULL PRIMARY KEY,
  createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  
  status orderStatus NOT NULL DEFAULT 'new',
  customerId int NOT NULL,
  FOREIGN KEY (customerId) REFERENCES customers(id)
);

CREATE TRIGGER setUpdatedAt
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE PROCEDURE triggerSetUpdatedAt();

CREATE TABLE IF NOT EXISTS order_products (
  id SERIAL NOT NULL PRIMARY KEY,
  createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  orderId int NOT NULL,
  productId int NOT NULL,
  unitPrice numeric(8,2) NOT NULL,
  quantity uint NOT NULL,

  FOREIGN KEY (orderId) REFERENCES orders(id),
  FOREIGN KEY (productId) REFERENCES products(id)
);

CREATE TRIGGER setUpdatedAt
BEFORE UPDATE ON order_products
FOR EACH ROW
EXECUTE PROCEDURE triggerSetUpdatedAt();
