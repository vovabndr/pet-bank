create table verify_emails
(
    id          bigserial primary key,
    username    varchar     not null,
    email       varchar     not null,
    secret_code varchar     not null,
    is_used     boolean     not null default false,
    created_at  timestamptz not null default (now()),
    expired_at  timestamptz not null default (now() + interval '15 minutes')
);

alter table verify_emails
    add foreign key ("username") references users ("username");

alter table users
    add column "is_email_verified" boolean not null default false;