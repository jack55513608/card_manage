-- +migrate Up
ALTER TABLE cards ADD COLUMN image_url TEXT;
