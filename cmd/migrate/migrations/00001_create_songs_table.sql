-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
  CREATE TABLE IF NOT EXISTS songs (
      id BIGSERIAL PRIMARY KEY,
      name TEXT NOT NULL,
      bpm INT NOT NULL CHECK (bpm BETWEEN 30 AND 246),
      sections JSONB NOT NULL,
       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );  


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd
