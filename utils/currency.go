package utils

const (
	USD = "USD"
	EUR = "EUR"
	UAH = "UAH"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, UAH:
		return true
	}
	return false
}
