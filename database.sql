/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE users (
	id serial PRIMARY KEY,
  phone VARCHAR(16) UNIQUE NOT NULL,
	name VARCHAR(64) NOT NULL,
  password VARCHAR NOT NULL,
  count_login INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (phone, name, password, count_login) VALUES ('+628123454321', 'Aldy Sanjaya Tanda', 'password', 5);
