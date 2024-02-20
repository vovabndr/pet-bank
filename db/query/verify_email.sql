-- name: CreateVerifyEmail :one
insert into verify_emails(username, email, secret_code)
values ($1, $2, $3)
returning *;

-- name: GetVerifyEmail :one
select *
from verify_emails
where id = $1
limit 1;

-- name: UpdateVerifyEmail :one
update verify_emails
set is_used = $1
where id = $2
returning *;