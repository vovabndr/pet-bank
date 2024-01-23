-- name: CreateTransfer :one
insert into transfers(from_account_id, to_account_id, amount) values ($1, $2, $3) returning *;

-- name: GetTransfer :one
select * from transfers where id = $1;