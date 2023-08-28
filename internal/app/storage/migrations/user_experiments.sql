CREATE TABLE user_experiments (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER,
    segment_id INTEGER REFERENCES Segments(id),
    UNIQUE(user_id, segment_id)
);