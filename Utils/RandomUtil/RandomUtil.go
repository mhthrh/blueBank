package RandomUtil

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdfghijkmnopuwqxyz"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOfList(list []string) string {
	return list[rand.Intn(len(list))]
}
