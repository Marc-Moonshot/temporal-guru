package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

func HashParams(params []string) string {
	paramsLength := len(params)
	if paramsLength == 0 {
		return ""
	}

	sorted := make([]string, paramsLength)
	copy(sorted, params)
	sort.Strings(sorted)

	joined := strings.Join(sorted, "&")

	hash := sha256.Sum256([]byte(joined))

	return hex.EncodeToString(hash[:])
}
