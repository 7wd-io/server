CREATE TABLE IF NOT EXISTS "game"
(
    id             serial PRIMARY KEY,
    host_nickname  varchar(15) NOT NULL,
    host_rating    int         NOT NULL,
    host_points    int         NOT NULL,
    guest_nickname varchar(15) NOT NULL,
    guest_rating   int         NOT NULL,
    guest_points   int         NOT NULL,
    winner         varchar(15) NULL,
    victory        smallint    NULL,
    log            jsonb       NOT NULL,
    started_at     TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    finished_at    TIMESTAMP(0) WITHOUT TIME ZONE NULL
);

CREATE INDEX IF NOT EXISTS ndx_game_hostnickname on game (host_nickname);
CREATE INDEX IF NOT EXISTS ndx_game_guestnickname on game (guest_nickname);
CREATE INDEX IF NOT EXISTS ndx_game_winner on game (winner);
CREATE INDEX IF NOT EXISTS ndx_game_victory on game (victory);
CREATE INDEX IF NOT EXISTS ndx_game_started_at on game (started_at);
