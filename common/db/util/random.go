package util

import (
	"fmt"
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

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomCustomerName() string {
	return RandomString(6)
}

func RandomCustomerAddr() string {
	return RandomString(10) + " " + RandomString(5)
}

func RandomPhone() string {
	return fmt.Sprintf("%d", RandomInt(1000000000, 9999999999))
}

func RandomEmail() string {
	return RandomString(10) + "@" + RandomString(5) + ".com"
}
