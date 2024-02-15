create table sessions(
    id uuid primary key,
    username varchar not null,
    refresh_token varchar not null,
    user_agent varchar not null,
    client_ip varchar not null,
    is_blocked boolean not null default false,
    expires_at timestamptz not null,
    created_at timestamptz not null default(now())
);

alter table sessions add foreign key ("username") references users("username");