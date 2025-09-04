package store

import "errors"

var (
	ErrAnalysisNotFound = errors.New("analysis not found")
	ErrEvidenceNotFound = errors.New("evidence not found")
)
