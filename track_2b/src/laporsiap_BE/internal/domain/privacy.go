package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

type PrivacyOptions struct {
	Mode  string   `json:"mode"`
	Terms []string `json:"terms,omitempty"`
}
type AIPreview struct {
	Revision   int64              `json:"revision"`
	Language   string             `json:"language"`
	Mode       string             `json:"mode"`
	Sources    []ChronologySource `json:"sources"`
	Redactions int                `json:"redactions"`
	Hash       string             `json:"hash"`
}

var emailPattern = regexp.MustCompile(`[A-Za-z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[A-Za-z0-9-]+(?:\.[A-Za-z0-9-]+)+`)
var phonePattern = regexp.MustCompile(`(?:\+62|62|0)8[0-9](?:[0-9 -]{6,17})[0-9]`)
var labelledNumber = regexp.MustCompile(`(?i)(?:\bNIK|\bKTP|\bidentity number|\bnational ID|\brekening|\baccount|\bacct)\s*[:#]?\s*([0-9]{6,20})\b`)

type redactionSpan struct {
	start, end int
	value      string
}

// Detection is deliberately limited. Preview and custom exact terms let the
// user check misses/false positives; this is not guaranteed anonymization.
func PrepareAIInput(c Case, options *PrivacyOptions, language string) (Case, AIPreview, error) {
	settings := PrivacyOptions{Mode: "redacted"}
	if options != nil {
		settings = *options
	}
	if settings.Mode != "redacted" && settings.Mode != "original" {
		return Case{}, AIPreview{}, errors.New("Pilih teks disamarkan atau teks asli untuk AI.")
	}
	if len(settings.Terms) > 50 {
		return Case{}, AIPreview{}, errors.New("Maksimal 50 istilah pribadi, masing-masing 2–160 karakter.")
	}
	terms := []string{}
	total := 0
	for _, term := range settings.Terms {
		term = strings.TrimSpace(term)
		n := utf8.RuneCountInString(term)
		total += len(term)
		if n < 2 || n > 160 || total > 4000 || !utf8.ValidString(term) {
			return Case{}, AIPreview{}, errors.New("Maksimal 50 istilah pribadi, masing-masing 2–160 karakter.")
		}
		terms = append(terms, term)
	}
	if settings.Mode == "redacted" {
		for _, m := range c.Messages {
			if len(strings.TrimSpace(m.Sender)) >= 3 {
				terms = append(terms, strings.TrimSpace(m.Sender))
			}
		}
	}
	prefix := "[LS_PRIVATE_"
	for strings.Contains(c.RawChat, prefix) {
		prefix = "[" + prefix
	}
	replacements := map[string]string{}
	termPatterns := []*regexp.Regexp{}
	seenTerms := map[string]bool{}
	for _, term := range terms {
		key := strings.ToLower(term)
		if !seenTerms[key] {
			seenTerms[key] = true
			termPatterns = append(termPatterns, regexp.MustCompile("(?i)"+regexp.QuoteMeta(term)))
		}
	}
	count := 0
	masked := Case{Messages: make([]Message, len(c.Messages))}
	for i, m := range c.Messages {
		text := m.Text
		if settings.Mode == "redacted" {
			spans := []redactionSpan{}
			for _, pattern := range []*regexp.Regexp{emailPattern, phonePattern} {
				for _, idx := range pattern.FindAllStringIndex(text, -1) {
					spans = append(spans, redactionSpan{idx[0], idx[1], text[idx[0]:idx[1]]})
				}
			}
			for _, idx := range labelledNumber.FindAllStringSubmatchIndex(text, -1) {
				spans = append(spans, redactionSpan{idx[2], idx[3], text[idx[2]:idx[3]]})
			}
			for _, pattern := range termPatterns {
				for _, idx := range pattern.FindAllStringIndex(text, -1) {
					spans = append(spans, redactionSpan{idx[0], idx[1], text[idx[0]:idx[1]]})
				}
			}
			sort.Slice(spans, func(a, b int) bool {
				if spans[a].start == spans[b].start {
					return spans[a].end > spans[b].end
				}
				return spans[a].start < spans[b].start
			})
			var result strings.Builder
			cursor := 0
			for _, span := range spans {
				if span.start < cursor {
					continue
				}
				result.WriteString(text[cursor:span.start])
				key := strings.ToLower(span.value)
				replacement, ok := replacements[key]
				if !ok {
					replacement = prefix + fmt.Sprintf("%04d", len(replacements)+1) + "]"
					replacements[key] = replacement
				}
				result.WriteString(replacement)
				cursor = span.end
				count++
			}
			result.WriteString(text[cursor:])
			text = result.String()
		}
		masked.Messages[i] = Message{ID: m.ID, Text: text}
		masked.RawChat += text + "\n"
	}
	preview := AIPreview{Revision: c.Revision, Language: language, Mode: settings.Mode, Sources: Sources(masked), Redactions: count}
	fingerprint, _ := json.Marshal(struct {
		CaseID  string
		Preview AIPreview
		Terms   []string
	}{c.ID, preview, settings.Terms})
	hash := sha256.Sum256(fingerprint)
	preview.Hash = hex.EncodeToString(hash[:])
	return masked, preview, nil
}
func RestoreEvidence(draft *Draft, masked, original Case, preview AIPreview) error {
	if draft == nil {
		return errors.New("Draft AI tidak memenuhi schema atau kutipan sumber. Draft ditolak; data asli tetap tersedia.")
	}
	raw, _ := json.Marshal(map[string]any{"events": draft.Events})
	events, err := ParseChronology(raw, Sources(masked))
	if err != nil {
		return &DraftValidationError{ResponseModel: draft.Model, Rule: err.Error()}
	}
	if err := ValidateGeneratedTiming(events); err != nil {
		return &DraftValidationError{ResponseModel: draft.Model, Rule: err.Error()}
	}
	for i := range events {
		events[i].Evidence = original.Messages[i].Text
	}
	draft.Events = events
	draft.Checks = ReviewChronology(events)
	draft.PrivacyMode = preview.Mode
	draft.Redactions = preview.Redactions
	return nil
}
