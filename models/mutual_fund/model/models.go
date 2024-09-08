// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package mutual_fund

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type MfInvestment struct {
	ID         int32
	SchemeID   pgtype.Int4
	Nav        pgtype.Numeric
	Units      pgtype.Numeric
	InvestedAt pgtype.Date
}

type MfNavDatum struct {
	SchemeID int32
	NavDate  pgtype.Date
	Nav      pgtype.Numeric
}
