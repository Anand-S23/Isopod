-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id            TEXT        PRIMARY KEY,
  github_id     BIGINT      NOT NULL,
  username      TEXT        NOT NULL,
  email         TEXT        NOT NULL DEFAULT '',
  avatar_url    TEXT        NOT NULL DEFAULT '',
  access_token  TEXT        NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL,
  last_login_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_users_email ON users (email);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS users;
