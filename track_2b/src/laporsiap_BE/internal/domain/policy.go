package domain

import "errors"

type SafetyLimits struct {
	RetentionDays     int   `json:"retention_days"`
	AccessHours       int   `json:"access_hours"`
	MaxCases          int   `json:"max_cases"`
	StorageBytes      int64 `json:"storage_bytes"`
	MaxRevisions      int64 `json:"max_revisions"`
	AICaseDaily       int   `json:"ai_case_daily"`
	AIGlobalDaily     int   `json:"ai_global_daily"`
	AIConcurrent      int   `json:"ai_concurrent"`
	RequestsPerMinute int   `json:"requests_per_minute"`
}

func DefaultSafetyLimits() SafetyLimits {
	return SafetyLimits{30, 168, 200, 256 * 1024 * 1024, 200, 3, 20, 2, 240}
}

var ErrCapacity = errors.New("Kapasitas penyimpanan atau revisi tercapai. Hapus kasus yang tidak diperlukan atau hubungi pengelola.")
var ErrAIQuota = errors.New("Batas panggilan AI harian tercapai. Draft manual tetap tersedia.")
var ErrAIBusy = errors.New("Apertus sedang melayani permintaan lain. Coba nanti atau gunakan draft manual.")

type ModelIdentityError struct {
	RequestedModel string `json:"requested_model"`
	ResponseModel  string `json:"response_model"`
}

func (*ModelIdentityError) Error() string {
	return "Identitas model respons bukan Apertus yang diminta. Draft ditolak; tidak ada fallback otomatis."
}

type DraftValidationError struct {
	ResponseModel string `json:"response_model"`
	Rule          string `json:"validation_rule"`
}

func (*DraftValidationError) Error() string {
	return "Draft AI tidak memenuhi schema atau kutipan sumber. Draft ditolak; data asli tetap tersedia."
}

var ErrNotFound = errors.New("Kasus tidak ditemukan atau akses tidak valid.")
var ErrConflict = errors.New("Kasus telah berubah. Muat ulang sebelum melanjutkan.")
