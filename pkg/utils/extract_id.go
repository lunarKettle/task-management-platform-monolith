package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func ExtractIDFromPath(path string) (uint32, error) {
	re := regexp.MustCompile(`^\/(?P<resource>[^\/]+(?:\/[^\/]+)*)\/(?P<id>\d+)$`)
	matches := re.FindStringSubmatch(path)
	if matches == nil {
		return 0, fmt.Errorf("path does not match expected format")
	}

	idIndex := re.SubexpIndex("id")
	if idIndex == -1 {
		return 0, fmt.Errorf("id group not found in path")
	}

	parsedUint, err := strconv.ParseUint(matches[idIndex], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid id format: %w", err)
	}

	return uint32(parsedUint), nil
}
