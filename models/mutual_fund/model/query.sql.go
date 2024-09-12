// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package mutual_fund

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addMFNavData = `-- name: AddMFNavData :exec
INSERT INTO mf_nav_data (scheme_id, nav_date, nav) VALUES ($1, $2, $3)
`

type AddMFNavDataParams struct {
	SchemeID int32
	NavDate  pgtype.Date
	Nav      pgtype.Numeric
}

func (q *Queries) AddMFNavData(ctx context.Context, arg AddMFNavDataParams) error {
	_, err := q.db.Exec(ctx, addMFNavData, arg.SchemeID, arg.NavDate, arg.Nav)
	return err
}

const addMFScheme = `-- name: AddMFScheme :exec
INSERT INTO mf_schemes (id, scheme_name) VALUES ($1, $2)
`

type AddMFSchemeParams struct {
	ID         int32
	SchemeName pgtype.Text
}

func (q *Queries) AddMFScheme(ctx context.Context, arg AddMFSchemeParams) error {
	_, err := q.db.Exec(ctx, addMFScheme, arg.ID, arg.SchemeName)
	return err
}

const cleanupMFNavDataBySchemeId = `-- name: CleanupMFNavDataBySchemeId :exec
DELETE FROM mf_nav_data WHERE scheme_id = $1
`

func (q *Queries) CleanupMFNavDataBySchemeId(ctx context.Context, schemeID int32) error {
	_, err := q.db.Exec(ctx, cleanupMFNavDataBySchemeId, schemeID)
	return err
}

const listMFInvestments = `-- name: ListMFInvestments :many
SELECT id, scheme_id, nav, units, invested_at FROM mf_investments
`

func (q *Queries) ListMFInvestments(ctx context.Context) ([]MfInvestment, error) {
	rows, err := q.db.Query(ctx, listMFInvestments)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MfInvestment
	for rows.Next() {
		var i MfInvestment
		if err := rows.Scan(
			&i.ID,
			&i.SchemeID,
			&i.Nav,
			&i.Units,
			&i.InvestedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMFInvestmentsBySchemeId = `-- name: ListMFInvestmentsBySchemeId :many
SELECT id, scheme_id, nav, units, invested_at FROM mf_investments WHERE scheme_id = $1
`

func (q *Queries) ListMFInvestmentsBySchemeId(ctx context.Context, schemeID pgtype.Int4) ([]MfInvestment, error) {
	rows, err := q.db.Query(ctx, listMFInvestmentsBySchemeId, schemeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MfInvestment
	for rows.Next() {
		var i MfInvestment
		if err := rows.Scan(
			&i.ID,
			&i.SchemeID,
			&i.Nav,
			&i.Units,
			&i.InvestedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMFNavData = `-- name: ListMFNavData :one
SELECT scheme_id, nav_date, nav FROM mf_nav_data WHERE scheme_id = $1 AND nav_date = $2
`

type ListMFNavDataParams struct {
	SchemeID int32
	NavDate  pgtype.Date
}

func (q *Queries) ListMFNavData(ctx context.Context, arg ListMFNavDataParams) (MfNavDatum, error) {
	row := q.db.QueryRow(ctx, listMFNavData, arg.SchemeID, arg.NavDate)
	var i MfNavDatum
	err := row.Scan(&i.SchemeID, &i.NavDate, &i.Nav)
	return i, err
}

const listMFNavDataBySchemeId = `-- name: ListMFNavDataBySchemeId :many
SELECT scheme_id, nav_date, nav FROM mf_nav_data WHERE scheme_id = $1 ORDER BY nav_date DESC
`

func (q *Queries) ListMFNavDataBySchemeId(ctx context.Context, schemeID int32) ([]MfNavDatum, error) {
	rows, err := q.db.Query(ctx, listMFNavDataBySchemeId, schemeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MfNavDatum
	for rows.Next() {
		var i MfNavDatum
		if err := rows.Scan(&i.SchemeID, &i.NavDate, &i.Nav); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMFSchemeById = `-- name: ListMFSchemeById :one
SELECT id, scheme_name FROM mf_schemes WHERE id = $1
`

func (q *Queries) ListMFSchemeById(ctx context.Context, id int32) (MfScheme, error) {
	row := q.db.QueryRow(ctx, listMFSchemeById, id)
	var i MfScheme
	err := row.Scan(&i.ID, &i.SchemeName)
	return i, err
}
