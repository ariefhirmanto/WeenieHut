-- +goose Up
-- +goose StatementBegin
CREATE TABLE activities (
    id BIGSERIAL PRIMARY KEY,                 
    user_id BIGINT NOT NULL,                 
    activity_type VARCHAR(20) NOT NULL,
    done_at TIMESTAMP NOT NULL,
    duration_minutes INT NOT NULL,
    calories_burned INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_activities_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_activities_user ON activities(user_id);

CREATE INDEX idx_activities_type ON activities(user_id, activity_type);

CREATE INDEX idx_activities_done_at ON activities(user_id, done_at);

CREATE INDEX idx_activities_calories ON activities(user_id, calories_burned);

CREATE INDEX idx_activities_search ON activities(user_id, activity_type, done_at, calories_burned);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_activities_search;
DROP INDEX IF EXISTS idx_activities_calories;
DROP INDEX IF EXISTS idx_activities_done_at;
DROP INDEX IF EXISTS idx_activities_type;
DROP INDEX IF EXISTS idx_activities_user;

DROP TABLE IF EXISTS activities;
-- +goose StatementEnd
