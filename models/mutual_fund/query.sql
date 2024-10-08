-- name: ListMFSchemeById :one
SELECT * FROM mf_schemes WHERE id = $1;

-- name: AddMFScheme :exec
INSERT INTO mf_schemes (id, scheme_name) VALUES ($1, $2);

-- name: ListDistinctMfInvestmentSchemeIds :many
SELECT DISTINCT scheme_id FROM mf_investments;

-- name: ListMFInvestmentsBySchemeId :many
SELECT * FROM mf_investments WHERE scheme_id = $1;

-- name: ListMFInvestments :many
SELECT * FROM mf_investments;

-- name: AddMFInvestment :exec
INSERT INTO mf_investments (scheme_id, nav, units, invested_at, created_at) VALUES ($1, $2, $3, $4, now());

-- name: AddMFNavData :exec
INSERT INTO mf_nav_data (scheme_id, nav_date, nav, created_at) VALUES ($1, $2, $3, now());

-- name: ListMFNavData :one
SELECT * FROM mf_nav_data WHERE scheme_id = $1 AND nav_date = $2;

-- name: ListMFNavDataBySchemeId :many
SELECT * FROM mf_nav_data WHERE scheme_id = $1 ORDER BY nav_date DESC;

-- name: CleanupMFNavDataBySchemeId :exec
DELETE FROM mf_nav_data WHERE scheme_id = $1;

