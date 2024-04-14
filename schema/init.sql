CREATE TABLE banner (
                         banner_id SERIAL PRIMARY KEY,
                         tag_ids INTEGER[] NOT NULL DEFAULT '{}',
                         feature_id INTEGER NOT NULL,
                         title VARCHAR(255) NOT NULL,
                         text TEXT NOT NULL,
                         url VARCHAR(255) NOT NULL,
                         is_active BOOLEAN NOT NULL DEFAULT FALSE,
                         created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users(
                      user_id SERIAL PRIMARY KEY,
                      email VARCHAR(255) NOT NULL UNIQUE,
                      password VARCHAR(255) NOT NULL,
                      is_admin BOOLEAN NOT NULL DEFAULT FALSE
);