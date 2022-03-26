CREATE TABLE IF NOT EXISTS events (
    id    bigserial   PRIMARY KEY,
    user_id     bigserial   NOT NULL,
    name       text        NOT NULL,
    details        text        NOT NULL,
    start   timestamp(0) with time zone NOT NULL,
    "end"  timestamp(0) with time zone NOT NULL,
    color text NOT NULL
);
