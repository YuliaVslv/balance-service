CREATE TABLE accounts (
    id serial NOT NULL PRIMARY KEY,
    balance FLOAT NOT NULL
);

CREATE TABLE transactions (
    id serial NOT NULL PRIMARY KEY,
    type VARCHAR(15) NOT NULL,
    user_id INTEGER NOT NULL,
    value FLOAT NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_DATE,
    description VARCHAR(255) NOT NULL
);

CREATE TABLE reserve (
    id serial NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    service_id INTEGER NOT NULL,
    order_id INTEGER NOT NULL,
    value FLOAT NOT NULL
);

CREATE TABLE reports
(
    year INTEGER NOT NULL DEFAULT EXTRACT(YEAR FROM CURRENT_DATE),
    month INTEGER NOT NULL DEFAULT EXTRACT(MONTH FROM CURRENT_DATE),
    service_id INTEGER NOT NULL,
    revenue FLOAT NOT NULL,
    PRIMARY KEY (year, month, service_id)
);
