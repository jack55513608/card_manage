-- +migrate Down
ALTER TABLE cards DROP COLUMN image_url;
