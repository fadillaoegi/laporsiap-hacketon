package config

import (
	"errors"
	"fmt"
	"laporsiap/internal/domain"
	"os"
	"strconv"
)

func LoadSafetyLimits() (domain.SafetyLimits, error) {
	limits := domain.DefaultSafetyLimits()
	for _, item := range []struct {
		name   string
		target *int
		max    int
	}{
		{"CASE_RETENTION_DAYS", &limits.RetentionDays, 365}, {"CASE_ACCESS_HOURS", &limits.AccessHours, 8760},
		{"MAX_CASES", &limits.MaxCases, 10000}, {"AI_CASE_DAILY_LIMIT", &limits.AICaseDaily, 100},
		{"AI_GLOBAL_DAILY_LIMIT", &limits.AIGlobalDaily, 1000}, {"AI_CONCURRENT_LIMIT", &limits.AIConcurrent, 10},
		{"REQUESTS_PER_MINUTE", &limits.RequestsPerMinute, 10000},
	} {
		if text := os.Getenv(item.name); text != "" {
			n, err := strconv.Atoi(text)
			if err != nil || n < 1 || n > item.max {
				return limits, fmt.Errorf("%s must be between 1 and %d", item.name, item.max)
			}
			*item.target = n
		}
	}
	for _, item := range []struct {
		name     string
		target   *int64
		min, max int64
	}{
		{"MAX_STORAGE_BYTES", &limits.StorageBytes, 1024 * 1024, 10 * 1024 * 1024 * 1024},
		{"MAX_CASE_REVISIONS", &limits.MaxRevisions, 10, 1000},
	} {
		if text := os.Getenv(item.name); text != "" {
			n, err := strconv.ParseInt(text, 10, 64)
			if err != nil || n < item.min || n > item.max {
				return limits, fmt.Errorf("%s is outside the supported range", item.name)
			}
			*item.target = n
		}
	}
	if limits.AccessHours > limits.RetentionDays*24 {
		return limits, errors.New("CASE_ACCESS_HOURS must not exceed case retention")
	}
	return limits, nil
}
