package utils

const (
	USD = "USD"
	INR = "INR"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, INR, CAD:
		return true
	}

	return false
}
