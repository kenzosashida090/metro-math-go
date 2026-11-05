package env

import "os"

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	print(key, "---FOUNDKEY---")
	if !ok {
		return fallback
	}
	return val
}
