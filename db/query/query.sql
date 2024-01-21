-- name: CreateAccount :one
insert into accounts(owner, balance, currency) values ($1, $2, $3) returning *;

-- name: GetAccount :one
select * from accounts where id = $1;

-- name: GetAccountForUpdate :one
select * from accounts where id = $1 for no key update;

-- name: ListAccounts :many
select * from accounts
order by id
limit $1
offset $2;

-- name: UpdateAccountBalance :one
update accounts set balance = $2 where id = $1 returning *;

-- name: AddAccountBalance :one
update accounts set balance = balance + sqlc.arg(amount) where id = sqlc.arg(id) returning *;

-- name: DeleteAccount :exec
delete from accounts where id = $1;


-- name: CreateTransfer :one
insert into transfers(from_account_id, to_account_id, amount) values ($1, $2, $3) returning *;

-- name: CreateEntry :one
insert into entries(account_id, amount) values ($1, $2) returning *;

-- name: GetTransfer :one
select * from transfers where id = $1;

-- name: GetEntry :one
select * from entries where id = $1;