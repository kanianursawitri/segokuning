create table if not exists friends_counter (
    user_id BIGSERIAL not null references users(id) PRIMARY KEY,
    friend_count BIGSERIAL not null,
    created_at timestamptz not null default current_timestamp
);

-- create index on user_id and friend_id
create index on friends(user_id);