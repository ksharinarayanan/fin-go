-- name: ListMFSchemeById :one
SELECT * FROM mf_schemes WHERE id = $1;

-- name: AddMFScheme :exec
INSERT INTO mf_schemes (id, scheme_name) VALUES ($1, $2);

-- name: ListMFInvestments :many
SELECT * FROM mf_investments;

-- name: AddMFNavData :exec
INSERT INTO mf_nav_data (scheme_id, nav_date, nav) VALUES ($1, $2, $3);

-- name: ListMFNavData :one
SELECT * FROM mf_nav_data WHERE scheme_id = $1 AND nav_date = $2;

-- name: ListMFNavDataBySchemeId :many
SELECT * FROM mf_nav_data WHERE scheme_id = $1 ORDER BY nav_date DESC;

-- name: CleanupMFNavDataBySchemeId :exec
DELETE FROM mf_nav_data WHERE scheme_id = $1;

