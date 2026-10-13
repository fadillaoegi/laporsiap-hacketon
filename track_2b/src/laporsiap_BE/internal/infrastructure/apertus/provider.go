package apertus

import (
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"laporsiap/internal/domain"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const ApertusModel = "swiss-ai/apertus-v1.5-70b"
const PromptVersion = "bilingual-v6-time-preserved"

// Only model identifiers cross the diagnostic boundary, never the raw response.
var modelIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9./_:-]{0,199}$`)

//go:embed prompt.json
var promptJSON []byte

// PromptSnapshot returns a copy of the exact embedded prompt used by this binary.
func PromptSnapshot() []byte { return append([]byte(nil), promptJSON...) }
func PromptSHA256() string {
	digest := sha256.Sum256(promptJSON)
	return hex.EncodeToString(digest[:])
}

type Provider struct {
	BaseURL, Key, Model string
	Client              *http.Client
}

func NewProvider(base, key, model string) *Provider {
	return &Provider{BaseURL: base, Key: key, Model: model, Client: &http.Client{Timeout: 40 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
}
func (p *Provider) Generate(ctx context.Context, c domain.Case, language string) (*domain.Draft, error) {
	started := time.Now()
	if language != "en" && language != "id" {
		return nil, errors.New("Bahasa draft harus Indonesia atau Inggris.")
	}
	if err := p.ValidateConfiguration(); err != nil {
		return nil, err
	}
	if len(c.Messages) > 20 || len(c.RawChat) > 12000 {
		return nil, errors.New("Draft AI dibatasi 20 pesan / 12 KB untuk menjaga kutipan lengkap. Gunakan draft manual untuk kasus lebih panjang.")
	}
	var prompt struct {
		System             string `json:"system_prompt"`
		Instruction        string `json:"user_instruction"`
		MaxTokens          int    `json:"max_tokens"`
		EnglishInstruction string `json:"english_instruction"`
	}
	_ = json.Unmarshal(promptJSON, &prompt)
	if language == "en" {
		prompt.Instruction = prompt.EnglishInstruction
	}
	msgs, _ := json.Marshal(domain.Sources(c))
	// No transfers, sender identities, attachment bytes, or case token leave this boundary.
	body, _ := json.Marshal(map[string]any{"model": p.Model, "temperature": 0, "max_tokens": prompt.MaxTokens, "messages": []map[string]string{{"role": "system", "content": prompt.System}, {"role": "user", "content": prompt.Instruction + "\n\nPesan sumber:\n" + string(msgs)}}})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint("chat/completions"), bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("Tidak dapat menyiapkan permintaan AI.")
	}
	req.Header.Set("Authorization", "Bearer "+p.Key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "LaporSiap/0.1")
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, errors.New("Apertus tidak merespons dalam batas waktu atau koneksi gagal. Coba lagi atau gunakan draft manual.")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Provider AI mengembalikan HTTP %d. Sumber tetap tersimpan; draft manual tersedia.", resp.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 256001))
	if err != nil || len(raw) > 256000 {
		return nil, errors.New("Respons AI terlalu besar atau gagal dibaca.")
	}
	var response struct {
		Model string `json:"model"`
		Usage *struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
	}
	if json.Unmarshal(raw, &response) != nil {
		return nil, errors.New("Respons provider tidak valid.")
	}
	if !strings.EqualFold(response.Model, p.Model) {
		modelID := "unrecognized"
		if modelIDPattern.MatchString(response.Model) {
			modelID = response.Model
		}
		return nil, &domain.ModelIdentityError{RequestedModel: p.Model, ResponseModel: modelID}
	}
	if len(response.Choices) != 1 || response.Choices[0].FinishReason != "stop" {
		return nil, errors.New("Respons AI belum lengkap. Gunakan draft manual atau coba kembali.")
	}
	events, err := domain.ParseChronologySummaries([]byte(response.Choices[0].Message.Content), domain.Sources(c))
	if err != nil {
		return nil, &domain.DraftValidationError{ResponseModel: response.Model, Rule: err.Error()}
	}
	draft := &domain.Draft{Origin: "apertus", Model: response.Model, Language: language, GeneratedAt: domain.Now(), DurationMS: time.Since(started).Milliseconds(), PromptVersion: PromptVersion, PromptSHA256: PromptSHA256(), Events: events, Checks: domain.ReviewChronology(events)}
	if response.Usage != nil && response.Usage.PromptTokens >= 0 && response.Usage.CompletionTokens >= 0 {
		draft.Usage = &domain.TokenUsage{Input: response.Usage.PromptTokens, Output: response.Usage.CompletionTokens}
	}
	return draft, nil
}
