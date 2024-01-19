package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "qwertyuioplkjhgfdsazxcvbnm"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(n int) string {
	var sb strings.Builder
	alphabetCount := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(alphabetCount)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return randomString(6)
}

func RandomMoney() int64 {
	return randomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "UAH"}
	l := len(currencies)
	return currencies[rand.Intn(l)]
}
