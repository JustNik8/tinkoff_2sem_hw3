CREATE TABLE IF NOT EXISTS "chat"
(
    "id"           UUID         NOT NULL,
    "nickname"     VARCHAR(255) NOT NULL,
    "message"      VARCHAR(255) NOT NULL,
    "message_time" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT current_timestamp
);
ALTER TABLE
    "chat"
    ADD PRIMARY KEY ("id");