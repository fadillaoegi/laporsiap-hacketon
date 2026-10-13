package postgres

import "laporsiap/internal/application"

var _ application.CaseRepository = (*Store)(nil)
