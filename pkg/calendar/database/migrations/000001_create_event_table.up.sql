CREATE TABLE IF NOT EXISTS events (
    event_id    bigserial   PRIMARY KEY,
    user_id     bigserial   NOT NULL,
    title       text        NOT NULL,
    body        text        NOT NULL,
    attend_at   timestamp(0) with time zone NOT NULL,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
