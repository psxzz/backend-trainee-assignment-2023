CREATE TABLE IF NOT EXISTS user_experiments (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER NOT NULL,
    segment_id INTEGER NOT NULL REFERENCES Segments(id),
    expires_at TIMESTAMP DEFAULT NULL,
    UNIQUE(user_id, segment_id)
);