package utils

import (
	"os"
	"strings"
)

// LoadEnv doing process split string from file .env
// and extract each key and value to os environment
func LoadEnv(env string) {
	s := strings.Split(env, "\r\n")

	for i, v := range s {
		if i < len(s) {
			// vS is the new slice of string that split by `=`,
			// the result is always 2 index, first index is key, and the second index is value
			vS := strings.Split(v, "=")

			os.Setenv(vS[0], vS[1])
		}

	}
}
