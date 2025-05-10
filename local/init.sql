SET TIME ZONE 'America/Puerto_Rico'

-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50),
    -- TAKEN FROM https://developer.spotify.com/documentation/web-api/reference/get-current-users-profile
    spotify_id VARCHAR(255) UNIQUE, -- USE TO ASSOCIATE SPOTIFY ID WITH OUR USER ID
    spotify_display_name VARCHAR(255),
    spotify_profile_uri VARCHAR(255),
    spotify_profile_picture_uri VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Songs Table
CREATE TABLE songs (
    song_id SERIAL PRIMARY KEY,
    spotify_id VARCHAR(255) UNIQUE,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255),
    album VARCHAR(255),
    release_date DATE,
    genre VARCHAR(100),
    cover_uri VARCHAR(255),
    preview_uri VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Friend Requests Table
CREATE TABLE friend_requests (
    request_id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    request_date TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50) CHECK (status IN ('pending', 'accepted', 'declined'))
);

-- Friends Table (Bi-Directional Friendships)
CREATE TABLE friends (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friendship_date DATE DEFAULT NOW(),
    PRIMARY KEY (user_id, friend_id)
);

-- Ensure Friendship Uniqueness (Avoid Duplicates Like (1,2) and (2,1))
CREATE UNIQUE INDEX unique_friendship 
ON friends (LEAST(user_id, friend_id), GREATEST(user_id, friend_id));


-- Rankings Table (Tracks Ratings for Songs by Users)
CREATE TABLE rankings (
    ranking_id SERIAL PRIMARY KEY,
    song_id INTEGER NOT NULL REFERENCES songs(song_id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rank INTEGER CHECK (rank >= 1 AND rank <= 5),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Index for Faster Ranking Queries
CREATE INDEX idx_rankings_user_song ON rankings(user_id, song_id);

-- Tracks JWT refresh tokens and rotations
CREATE TABLE jwt_refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    jti TEXT NOT NULL, -- new column for the token identifier
    refresh_token TEXT NOT NULL, 
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE spotify_refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL, 
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

create table streaks(
    user_id          integer                 not null
        primary key
        references users,
    streak_count     integer   default 0     not null,
    daily_count      integer   default 0     not null,
    last_count_date  date,
    last_streak_date date,
    updated_at       timestamp default now() not null
);


CREATE TABLE impression_stats (
  impression_label TEXT PRIMARY KEY UNIQUE NOT NULL,
  impressions      BIGINT             NOT NULL DEFAULT 0,
  clicks           BIGINT             NOT NULL DEFAULT 0,
  created_at       TIMESTAMP NOT NULL DEFAULT NOW()
);

--give ownership to ranktifyUser
ALTER TABLE users OWNER TO ranktifyUser;
ALTER TABLE songs OWNER TO ranktifyUser;
ALTER TABLE friend_requests OWNER TO ranktifyUser;
ALTER TABLE friends OWNER TO ranktifyUser;
ALTER TABLE rankings OWNER TO ranktifyUser;
ALTER TABLE jwt_refresh_tokens OWNER TO ranktifyUser;
ALTER TABLE spotify_refresh_tokens OWNER TO ranktifyUser;
ALTER TABLE impression_stats OWNER TO ranktifyUser;