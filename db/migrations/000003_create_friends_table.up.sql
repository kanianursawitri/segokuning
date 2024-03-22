create table if not exists friends(
    id serial primary key,
    user_id int not null references users(id),
    friend_id int not null references users(id),
    created_at timestamptz not null default current_timestamp
);

-- create index on user_id and friend_id
create index on friends(user_id);
create index on friends(friend_id);
-- create unique index on user_id and friend_id
create unique index on friends(user_id, friend_id);