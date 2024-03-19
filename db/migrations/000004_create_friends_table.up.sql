create table if not exists friends(
    first_user_id bigserial not null,
    second_user_id bigserial not null,
    primary key (first_user_id, second_user_id),
    created_at timestamptz not null default current_timestamp
);

-- Create indexes
create index first_user_id on friends (first_user_id);
create index second_user_id on friends (second_user_id);