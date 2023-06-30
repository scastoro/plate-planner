-- +goose Up
CREATE SCHEMA "Admin";

CREATE TABLE "Admin"."Users" (
  id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  body_weight DECIMAL NOT NULL,
  username TEXT NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  lastLoggedIn TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE "Admin"."Users";

DROP SCHEMA "Admin";