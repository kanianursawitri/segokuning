create table if not exists friends(
    first_user_id bigserial not null,
    second_user_id bigserial not null,
    primary key (first_user_id, second_user_id)
    created_at timestamptz not null default current_timestamp
);