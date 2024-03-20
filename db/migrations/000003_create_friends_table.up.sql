create table if not exists friends(
    id serial primary key,
    user_id int not null references users(id),
    friend_id int not null references users(id),
    created_at timestamptz not null default current_timestamp
);
