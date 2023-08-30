CREATE TYPE user_experiments_op AS ENUM('add', 'remove');
CREATE TABLE IF NOT EXISTS log_user_experiments (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER NOT NULL,
    segment_name VARCHAR(256) NOT NULL,
    op_type user_experiments_op NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW()
);