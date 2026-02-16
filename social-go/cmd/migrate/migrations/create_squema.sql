CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email citext NOT NULL UNIQUE,
    password bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts(
                                    id BIGSERIAL PRIMARY KEY,
                                    title VARCHAR(255),
                                    user_id BIGINT REFERENCES users(id),
                                    content TEXT NOT NULL,
                                    version int NOT NULL DEFAULT 0,
                                    tags TEXT[],
                                    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
                                    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments(
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT REFERENCES posts(id),
    user_id BIGINT REFERENCES users(id),
    content TEXT NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS followers (
     id BIGSERIAL PRIMARY KEY,
     user_id BIGINT NOT NULL REFERENCES users(id),
     follower_id BIGINT NOT NULL REFERENCES users(id),
     created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
     CONSTRAINT followers_user_follower_unique UNIQUE (user_id, follower_id)
);
