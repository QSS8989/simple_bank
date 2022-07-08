package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		b := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(b)
	}

	return sb.String()
}

func RandOwner() string {
	return RandString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	return currencies[rand.Intn(len(currencies))]
}
