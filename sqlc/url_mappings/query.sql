-- name: CreateUrlMapping :one
insert into url_mappings (shortened_url, url, inactive_expire_at, must_expire_at)
values ($1, $2, $3, $4)
returning *;

-- name: GetUrlMapping :one
select *
from url_mappings
where shortened_url = $1
limit 1;

-- name: DeleteUrlMappingWithId :exec
delete from url_mappings
where id = $1;

-- name: DeleteExpiredUrlMappings :exec
delete from url_mappings
where timezone('UTC', now()) >= url_mappings.inactive_expire_at
   or timezone('UTC', now()) >= must_expire_at;