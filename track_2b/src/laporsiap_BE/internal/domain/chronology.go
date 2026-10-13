package domain

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strings"
)

// These conservative checks flag observed failure modes for human review.
// They do not prove semantic entailment or automatically rewrite a draft.
type ChronologySource struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type ChronologyEvent struct {
	SourceID    string `json:"source_id"`
	Description string `json:"description"`
	Evidence    string `json:"evidence"`
}
type ChronologyReview struct {
	EventIndex int      `json:"event_index"`
	SourceID   string   `json:"source_id"`
	Flags      []string `json:"flags"`
}

var sentWords = regexp.MustCompile(`\b(mengirim|kirim|mentransfer|membayar|bayar|menyetor|setor|topup|top-up|send|sent|sending|pay|paid|paying|transferred)\b`)
var receivedWords = regexp.MustCompile(`\b(menerima|diterima|masuk|dikreditkan|receive|received|receiving|credited)\b`)
var financialWords = regexp.MustCompile(`\b(rp|idr|usd|chf)\s*\d|\b(uang|dana|transfer|transferred|money|funds|payment|deposit|topup|top-up|setor|setoran)\b`)
var instructionWords = regexp.MustCompile(`\bsystem\s*:|\bsystem override\b|\binstruksi untuk ai\b|\bperintah untuk ai\b|\b(ignore|disregard) (all |the |previous |prior )*instructions\b|\b(abaikan|lupakan) (semua |seluruh |instruksi )*(instruksi|perintah)\b|\binstructions? (for|to) (the )?ai\b`)
var uncertainTransfer = regexp.MustCompile(`\b(pending|failed|gagal|dalam proses|belum berhasil|belum selesai|belum membayar|belum mengirim|not completed|not successful|not sent|not paid|unsuccessful|did not (send|pay|transfer)|have not (sent|paid|transferred)|has not (sent|paid|transferred)|balance (did not change|unchanged))\b`)
var completedTransfer = regexp.MustCompile(`\b(berhasil|selesai|sudah (membayar|mengirim|ditransfer)|successful|successfully|completed|has paid|has sent|payment went through)\b`)
var negatedCompletion = regexp.MustCompile(`\b(belum|tidak|not|never|no|unsuccessful|failed|gagal)\b`)
var requestOrPromise = regexp.MustCompile(`\b(meminta|diminta|minta|silakan|tolong|menjanjikan|janji|dijanjikan|akan|jika|agar|asks?|asked|requests?|requested|promises?|promised|will|would|please|if|to unlock|in order to)\b`)
var attributedClaim = regexp.MustCompile(`\b(meminta|diminta|minta|menjanjikan|janji|dijanjikan|mengklaim|menyatakan|mengatakan|klaim|akan|jika|agar|belum|tidak|asks?|asked|requests?|requested|promises?|promised|claims?|claimed|states?|stated|says?|said|reports?|reported|will|would|if|not)\b`)
var rupiahAmounts = regexp.MustCompile(`\b(?:rp|idr)\s*(\d+(?:[.,]\d{3})*)\b`)
var recruiterRole = regexp.MustCompile(`\b(perekrut|recruiter)\b`)

// These explicit relative-time phrases are compared across EN/ID. Generated
// drafts must preserve them; user edits keep review flags. This does not infer
// calendar dates or establish full semantic entailment.
var relativeTimes = []*regexp.Regexp{
	regexp.MustCompile(`\b(tomorrow|besok)\b`),
	regexp.MustCompile(`\b(next week|pekan depan|minggu depan)\b`),
	regexp.MustCompile(`\b(next month|bulan depan)\b`),
	regexp.MustCompile(`\b(today|hari ini)\b`),
	regexp.MustCompile(`\b(yesterday|kemarin)\b`),
}

func amountSet(text string) map[string]bool {
	out := map[string]bool{}
	for _, match := range rupiahAmounts.FindAllStringSubmatch(text, -1) {
		n := strings.NewReplacer(".", "", ",", "").Replace(match[1])
		out[strings.TrimLeft(n, "0")] = true
	}
	return out
}

func ParseChronology(raw []byte, sources []ChronologySource) ([]ChronologyEvent, error) {
	var draft struct {
		Events []ChronologyEvent `json:"events"`
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&draft); err != nil {
		return nil, errors.New("invalid chronology JSON schema")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, errors.New("trailing JSON after chronology")
	}
	if draft.Events == nil || len(draft.Events) != len(sources) {
		return nil, errors.New("chronology event count differs from sources")
	}
	seen := map[string]bool{}
	for i, event := range draft.Events {
		if seen[event.SourceID] || event.SourceID != sources[i].ID {
			return nil, errors.New("chronology source order or identifier invalid")
		}
		if event.Evidence != sources[i].Text {
			return nil, errors.New("chronology quote differs from original source")
		}
		if strings.TrimSpace(event.Description) == "" || len(strings.Fields(event.Description)) > 35 {
			return nil, errors.New("chronology description empty or over 35 words")
		}
		seen[event.SourceID] = true
	}
	return draft.Events, nil
}

// ParseChronologySummaries accepts only source IDs and descriptions from AI.
// Evidence is bound locally after checking count/order/IDs; generated quotes
// are never accepted or repaired. ParseChronology remains the full event
// contract for saved drafts, user edits, exports, and privacy restoration.
func ParseChronologySummaries(raw []byte, sources []ChronologySource) ([]ChronologyEvent, error) {
	var response struct {
		Events []struct {
			SourceID    string `json:"source_id"`
			Description string `json:"description"`
		} `json:"events"`
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&response); err != nil {
		return nil, errors.New("invalid chronology JSON schema")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, errors.New("trailing JSON after chronology")
	}
	if response.Events == nil || len(response.Events) != len(sources) {
		return nil, errors.New("chronology event count differs from sources")
	}
	events := make([]ChronologyEvent, len(sources))
	for i, event := range response.Events {
		events[i] = ChronologyEvent{SourceID: event.SourceID, Description: event.Description, Evidence: sources[i].Text}
	}
	bound, _ := json.Marshal(map[string]any{"events": events})
	validated, err := ParseChronology(bound, sources)
	if err != nil {
		return nil, err
	}
	if err := ValidateGeneratedTiming(validated); err != nil {
		return nil, err
	}
	return validated, nil
}

func relativeTimingDiffers(source, summary string) bool {
	source, summary = strings.ToLower(source), strings.ToLower(summary)
	for _, expression := range relativeTimes {
		if expression.MatchString(source) != expression.MatchString(summary) {
			return true
		}
	}
	return false
}

// ValidateGeneratedTiming rejects omissions, changes, and invented phrases
// before saving an AI draft. It never appends a date or rewrites a summary.
func ValidateGeneratedTiming(events []ChronologyEvent) error {
	for _, event := range events {
		if relativeTimingDiffers(event.Evidence, event.Description) {
			return errors.New("chronology relative time differs from source")
		}
	}
	return nil
}

func ReviewChronology(events []ChronologyEvent) []ChronologyReview {
	checks := make([]ChronologyReview, 0, len(events))
	for i, event := range events {
		quote := strings.ToLower(event.Evidence)
		summary := strings.ToLower(event.Description)
		flags := []string{}
		if instructionWords.MatchString(quote) {
			flags = append(flags, "instruction_in_source_requires_review")
		}
		if relativeTimingDiffers(quote, summary) {
			flags = append(flags, "relative_time_requires_review")
		}
		quoteAmounts := amountSet(quote)
		summaryAmounts := amountSet(summary)
		for amount := range quoteAmounts {
			if !summaryAmounts[amount] {
				flags = append(flags, "source_amount_omitted_requires_review")
				break
			}
		}
		for amount := range summaryAmounts {
			if !quoteAmounts[amount] {
				flags = append(flags, "amount_outside_source_requires_review")
				break
			}
		}
		if financialWords.MatchString(quote) {
			quoteSent, quoteReceived := sentWords.MatchString(quote), receivedWords.MatchString(quote)
			summarySent, summaryReceived := sentWords.MatchString(summary), receivedWords.MatchString(summary)
			if quoteSent && !quoteReceived && summaryReceived && !summarySent ||
				quoteReceived && !quoteSent && summarySent && !summaryReceived {
				flags = append(flags, "financial_direction_requires_review")
			}
			if quoteSent && !summarySent && !requestOrPromise.MatchString(quote) {
				flags = append(flags, "financial_action_omitted_requires_review")
			}
			if recruiterRole.MatchString(quote) && requestOrPromise.MatchString(quote) && !recruiterRole.MatchString(summary) {
				flags = append(flags, "quoted_role_requires_review")
			}
			if uncertainTransfer.MatchString(quote) && completedTransfer.MatchString(summary) && !negatedCompletion.MatchString(summary) {
				flags = append(flags, "transfer_status_requires_review")
			}
			if requestOrPromise.MatchString(quote) && !attributedClaim.MatchString(summary) && (sentWords.MatchString(summary) || receivedWords.MatchString(summary) || completedTransfer.MatchString(summary)) {
				flags = append(flags, "claim_as_fact_requires_review")
			}
		}
		// Only a small documented vocabulary is checked. Legitimate paraphrases
		// may also be flagged; the user's review decides, not this heuristic.
		for _, term := range []string{"upgrade", "member", "selesai"} {
			if strings.Contains(summary, term) && !strings.Contains(quote, term) {
				flags = append(flags, "summary_term_outside_quote:"+term)
			}
		}
		if len(flags) > 0 {
			checks = append(checks, ChronologyReview{EventIndex: i, SourceID: event.SourceID, Flags: flags})
		}
	}
	return checks
}
