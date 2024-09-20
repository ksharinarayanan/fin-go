package mutual_fund

import "errors"

var ErrInvalidSchemeId = errors.New("invalid scheme id")

var ErrNoInvestmentsForSchemeId = errors.New("no investments for scheme")
