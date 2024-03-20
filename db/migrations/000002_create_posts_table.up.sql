create table if not exists posts(
    id bigserial primary key,
    user_id bigserial references users(id) on delete cascade,
    post_in_html varchar not null,
    tags varchar[] not null default array[]::varchar[],
    comments jsonb not null default '[]'::jsonb,
    created_at timestamptz not null default current_timestamp
);

-- Create indexes
create index user_id on posts (user_id);