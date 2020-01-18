package utils

import "math/rand"

const digitBytes = "123456789"
const (
	digitIdxBits = 6
	digitIdxMask = 1<<digitIdxBits - 1
	digitIdxMax  = 63 / digitIdxBits
)

func GetRandDigitsString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), digitIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), digitIdxMax
		}
		if idx := int(cache & digitIdxMask); idx < len(digitBytes) {
			b[i] = digitBytes[idx]
			i--
		}
		cache >>= digitIdxBits
		remain--
	}
	return string(b)
}
