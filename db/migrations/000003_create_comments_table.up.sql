create table if not exists comments(
    id bigserial primary key,
    post_id bigserial references posts(id) on delete cascade,
    user_id bigserial references users(id) on delete cascade,
    comment varchar not null,
    created_at timestamptz not null default current_timestamp
);

-- Create indexes
create index post_id on comments (post_id);