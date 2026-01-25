package efir

import "math/rand/v2"

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const ln = 8

func addSuffix(in string) string {
	alphaLn := len(alphabet)
	suffix := make([]byte, ln)

	for i := range suffix {
		suffix[i] = alphabet[rand.IntN(alphaLn)]
	}

	return in + "-" + string(suffix)
}
