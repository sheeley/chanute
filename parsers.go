package chanute

import (
	"fmt"
	"strconv"
	"strings"
)

func parseAmount(s string) int {
	s = strings.ReplaceAll(s, "$", "")
	idx := strings.Index(s, ".")
	if idx != -1 {
		s = s[0:idx]
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return i
}

func parseDays(s string) int {
	if s == "14+" {
		return 15
	}
	s = strings.TrimSpace(strings.ReplaceAll(s, "days", ""))
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return i
}
