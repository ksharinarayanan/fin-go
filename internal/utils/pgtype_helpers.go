package utils

import (
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgNumericToFloat64(n pgtype.Numeric) (float64, error) {
	value, err := n.Float64Value()
	return value.Float64, err
}

func TimeToPgDate(date time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  date,
		Valid: true,
	}
}

func Float64ToPgNumeric(val float64) pgtype.Numeric {
	parse := strconv.FormatFloat(val, 'f', -1, 64)
	var result pgtype.Numeric
	if err := result.Scan(parse); err != nil {
		log.Println("Error scanning numeric")
	}
	return result
}

func StringToPgText(val string) pgtype.Text {
	return pgtype.Text{
		String: val,
		Valid:  true,
	}
}

func IntToPgInt4(val int) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(val),
		Valid: true,
	}
}
