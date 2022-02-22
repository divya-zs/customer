--DROP DATABASE IF EXISTS customers;

--CREATE DATABASE customers;
-- \c customers;
CREATE SCHEMA IF NOT EXISTS test AUTHORIZATION postgres;
DROP TABLE IF EXISTS customers;

CREATE TABLE customer(
                    id SERIAL PRIMARY KEY,
                    name varchar(20) NOT NULL UNIQUE,
                    age int,
                    salary float
);

INSERT INTO customer VALUES(1,'Divya', 22, 30000);
INSERT INTO customer VALUES(2,'Jay', 21, 30000);
INSERT INTO customer VALUES(3,'Karan', 22, 30000);

CREATE TABLE DELETED_USER(
    id int,
    name varchar(20),
    age int,
    salary float
);
--
-- CREATE FUNCTION moveDeleted() RETURNS trigger AS $$
--     BEGIN
--         INSERT INTO DELETED_USER VALUES(OLD);
--         RETURN OLD;
--     END;
-- $$ LANGUAGE plpgsql;