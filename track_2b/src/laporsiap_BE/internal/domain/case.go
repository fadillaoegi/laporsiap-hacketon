package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type Message struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Sender string `json:"sender"`
	SentAt string `json:"sent_at"`
}
type Transfer struct {
	ID          string `json:"id"`
	FromBank    string `json:"from_bank"`
	FromAccount string `json:"from_account"`
	ToBank      string `json:"to_bank"`
	ToAccount   string `json:"to_account"`
	Amount      *int64 `json:"amount"`
	OccurredAt  string `json:"occurred_at"`
	Status      string `json:"status"`
	Note        string `json:"note"`
}
type Attachment struct {
	ID          string   `json:"id"`
	Filename    string   `json:"filename"`
	ContentType string   `json:"content_type"`
	Size        int      `json:"size"`
	SHA256      string   `json:"sha256"`
	SourceIDs   []string `json:"source_ids,omitempty"`
	TransferIDs []string `json:"transfer_ids,omitempty"`
}
type TokenUsage struct {
	Input  int `json:"input_tokens"`
	Output int `json:"output_tokens"`
}
type Draft struct {
	Origin        string             `json:"origin"`
	Model         string             `json:"model"`
	Events        []ChronologyEvent  `json:"events"`
	Checks        []ChronologyReview `json:"checks"`
	Language      string             `json:"language,omitempty"`
	GeneratedAt   string             `json:"generated_at,omitempty"`
	DurationMS    int64              `json:"duration_ms,omitempty"`
	Usage         *TokenUsage        `json:"usage,omitempty"`
	PromptVersion string             `json:"prompt_version,omitempty"`
	PromptSHA256  string             `json:"prompt_sha256,omitempty"`
	PrivacyMode   string             `json:"privacy_mode,omitempty"`
	Redactions    int                `json:"redactions,omitempty"`
}
type Case struct {
	ID               string       `json:"id"`
	Title            string       `json:"title"`
	Synthetic        bool         `json:"synthetic"`
	RawChat          string       `json:"raw_chat"`
	Parser           string       `json:"parser"`
	Messages         []Message    `json:"messages"`
	Transfers        []Transfer   `json:"transfers"`
	Draft            *Draft       `json:"draft"`
	Attachments      []Attachment `json:"attachments"`
	Revision         int64        `json:"revision"`
	ReviewedRevision *int64       `json:"reviewed_revision"`
	ReviewedAt       string       `json:"reviewed_at"`
	CreatedAt        string       `json:"created_at"`
	UpdatedAt        string       `json:"updated_at"`
	ExpiresAt        string       `json:"expires_at"`
	AccessExpiresAt  string       `json:"access_expires_at"`
}
type Input struct {
	Title     string     `json:"title"`
	Synthetic bool       `json:"synthetic"`
	RawChat   string     `json:"raw_chat"`
	Parser    string     `json:"parser"`
	Transfers []Transfer `json:"transfers"`
}
type Summary struct {
	SuccessfulKnownTotal    *int64         `json:"successful_known_total"`
	SuccessfulCount         int            `json:"successful_count"`
	SuccessfulUnknownAmount int            `json:"successful_unknown_amount"`
	StatusCounts            map[string]int `json:"status_counts"`
	Warnings                []string       `json:"warnings"`
}

func RandomID() string {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
func TokenHash(token string) []byte { h := sha256.Sum256([]byte(token)); return h[:] }
func Now() string                   { return time.Now().UTC().Format(time.RFC3339) }

// Explicitly supported format: [DD/MM/YYYY, HH:MM(:SS)] Sender: message.
// Export wall-clock strings are retained verbatim; no timezone is invented.
var waHeader = regexp.MustCompile(`^\[(\d{2}/\d{2}/\d{4}, \d{2}:\d{2}(?::\d{2})?)\] ([^:]+): (.*)$`)

func ParseMessages(raw, format string) ([]Message, error) {
	raw = strings.ReplaceAll(strings.TrimPrefix(raw, "\ufeff"), "\r\n", "\n")
	if !utf8.ValidString(raw) || len(raw) > 32000 || strings.TrimSpace(raw) == "" {
		return nil, errors.New("Isi percakapan UTF-8 wajib diisi (maksimal 32 KB).")
	}
	out := []Message{}
	switch format {
	case "plain":
		for _, line := range strings.Split(raw, "\n") {
			if strings.TrimSpace(line) != "" {
				out = append(out, Message{Text: line})
			}
		}
	case "whatsapp_bracket_dmy":
		for _, line := range strings.Split(raw, "\n") {
			match := waHeader.FindStringSubmatch(line)
			if match != nil {
				layout := "02/01/2006, 15:04"
				if len(match[1]) > 17 {
					layout = "02/01/2006, 15:04:05"
				}
				if _, err := time.Parse(layout, match[1]); err != nil {
					return nil, errors.New("Tanggal/jam header TXT tidak valid. Pilih teks biasa untuk mempertahankan sumber tanpa metadata.")
				}
				out = append(out, Message{Text: match[3], Sender: strings.TrimSpace(match[2]), SentAt: match[1]})
			} else if len(out) == 0 {
				return nil, errors.New("Format TXT belum didukung. Gunakan [DD/MM/YYYY, HH:MM] Pengirim: pesan, atau pilih teks biasa.")
			} else if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
				return nil, errors.New("Header TXT ambigu. Pilih teks biasa agar metadata tidak ditebak.")
			} else {
				out[len(out)-1].Text += "\n" + line
			}
		}
	default:
		return nil, errors.New("Pilih format teks yang didukung.")
	}
	if len(out) == 0 || len(out) > 100 {
		return nil, errors.New("Gunakan 1–100 pesan per kasus.")
	}
	for i := range out {
		if strings.TrimSpace(out[i].Text) == "" || len(out[i].Text) > 6000 {
			return nil, errors.New("Setiap pesan harus berisi teks (maksimal 6 KB).")
		}
		out[i].ID = fmt.Sprintf("M%d", i+1)
	}
	return out, nil
}
func ValidateInput(in Input) ([]Message, error) {
	if strings.TrimSpace(in.Title) == "" || len(in.Title) > 160 {
		return nil, errors.New("Judul wajib diisi, maksimal 160 karakter.")
	}
	if len(in.Transfers) > 100 {
		return nil, errors.New("Maksimal 100 catatan transfer.")
	}
	seen := map[string]bool{}
	var total int64
	for i, t := range in.Transfers {
		if t.ID == "" || len(t.ID) > 100 || seen[t.ID] {
			return nil, errors.New("ID transfer harus unik dan tidak kosong.")
		}
		seen[t.ID] = true
		if t.Status != "successful" && t.Status != "pending" && t.Status != "failed" && t.Status != "unknown" {
			return nil, errors.New("Status transfer tidak valid.")
		}
		if t.Amount != nil {
			if *t.Amount < 0 || *t.Amount > 9000000000000000 {
				return nil, errors.New("Nominal harus berupa rupiah bulat nonnegatif, atau kosong.")
			}
			if t.Status == "successful" {
				if total > math.MaxInt64-*t.Amount {
					return nil, errors.New("Total transfer terlalu besar.")
				}
				total += *t.Amount
			}
		}
		if t.OccurredAt != "" {
			if _, err := time.Parse("2006-01-02", t.OccurredAt); err != nil {
				return nil, fmt.Errorf("Tanggal transfer %d harus YYYY-MM-DD atau kosong.", i+1)
			}
		}
		for _, s := range []string{t.FromBank, t.FromAccount, t.ToBank, t.ToAccount, t.Note} {
			if len(s) > 2000 {
				return nil, errors.New("Detail transfer terlalu panjang.")
			}
		}
	}
	if total > 9007199254740991 {
		return nil, errors.New("Total melebihi batas angka yang dapat ditampilkan secara presisi.")
	}
	return ParseMessages(in.RawChat, in.Parser)
}
func Summarize(c Case) Summary {
	s := Summary{StatusCounts: map[string]int{"successful": 0, "pending": 0, "failed": 0, "unknown": 0}, Warnings: []string{}}
	var total int64
	known := false
	pairs := map[string]bool{}
	for i, t := range c.Transfers {
		s.StatusCounts[t.Status]++
		if t.Status == "successful" {
			s.SuccessfulCount++
			if t.Amount == nil {
				s.SuccessfulUnknownAmount++
			} else {
				known = true
				total += *t.Amount
			}
		}
		fields := []string{}
		if t.Amount == nil {
			fields = append(fields, "nominal")
		}
		if t.OccurredAt == "" {
			fields = append(fields, "tanggal transfer")
		}
		if t.FromBank == "" || t.FromAccount == "" {
			fields = append(fields, "rekening pengirim")
		}
		if t.ToBank == "" || t.ToAccount == "" {
			fields = append(fields, "rekening penerima")
		}
		if t.Status == "unknown" {
			fields = append(fields, "status")
		}
		if len(fields) > 0 {
			s.Warnings = append(s.Warnings, fmt.Sprintf("Transfer %d: %s belum lengkap.", i+1, strings.Join(fields, ", ")))
		}
		key := fmt.Sprintf("%q/%q/%q/%q/%v/%q/%q", t.FromBank, t.FromAccount, t.ToBank, t.ToAccount, amountText(t.Amount), t.OccurredAt, t.Status)
		if pairs[key] {
			s.Warnings = append(s.Warnings, fmt.Sprintf("Transfer %d memiliki rincian yang sama dengan catatan lain; periksa duplikasi.", i+1))
		}
		pairs[key] = true
	}
	if known {
		s.SuccessfulKnownTotal = &total
	}
	if s.SuccessfulUnknownAmount > 0 {
		s.Warnings = append(s.Warnings, "Total sebagian: ada transfer berhasil dengan nominal belum diketahui.")
	}
	if len(c.Attachments) == 0 {
		s.Warnings = append(s.Warnings, "Belum ada lampiran bukti. Lampiran dapat ditambahkan tanpa OCR.")
	}
	return s
}
func Sources(c Case) []ChronologySource {
	s := make([]ChronologySource, len(c.Messages))
	for i, m := range c.Messages {
		s[i] = ChronologySource{ID: m.ID, Text: m.Text}
	}
	return s
}
func ValidateEditedEvents(events []ChronologyEvent, c Case) error {
	if len(events) != len(c.Messages) {
		return errors.New("Setiap pesan sumber wajib memiliki satu entri kronologi.")
	}
	for i, e := range events {
		if e.SourceID != c.Messages[i].ID || e.Evidence != c.Messages[i].Text || strings.TrimSpace(e.Description) == "" || len(e.Description) > 8000 {
			return errors.New("Sumber dan kutipan tidak dapat diubah melalui draft; ringkasan wajib diisi.")
		}
	}
	return nil
}
func ManualDraft(c Case) *Draft {
	return ManualDraftInLanguage(c, "id")
}
func ManualDraftInLanguage(c Case, language string) *Draft {
	prefix := "Kutipan pesan: "
	if language == "en" {
		prefix = "Message quote: "
	}
	events := make([]ChronologyEvent, len(c.Messages))
	for i, m := range c.Messages {
		events[i] = ChronologyEvent{SourceID: m.ID, Description: prefix + m.Text, Evidence: m.Text}
	}
	return &Draft{Origin: "manual", Language: language, GeneratedAt: Now(), Events: events, Checks: ReviewChronology(events)}
}
func Invalidate(c *Case) { c.ReviewedRevision = nil; c.ReviewedAt = "" }

func amountText(v *int64) string {
	if v == nil {
		return "Belum diketahui"
	}
	return fmt.Sprintf("%d", *v)
}
