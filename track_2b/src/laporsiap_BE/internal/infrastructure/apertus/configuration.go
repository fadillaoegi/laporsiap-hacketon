package apertus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ConfigFromEnv reads the Hack Apertus variables (LLM_BASE_URL, LLM_API_KEY,
// LLM_NAME) first and falls back to the earlier APERTUS_* names.
func ConfigFromEnv() (base, key, model string) {
	pick := func(names ...string) string {
		for _, name := range names {
			if value := strings.TrimSpace(os.Getenv(name)); value != "" {
				return value
			}
		}
		return ""
	}
	base = pick("LLM_BASE_URL", "APERTUS_API_BASE_URL")
	key = pick("LLM_API_KEY", "APERTUS_API_KEY")
	model = pick("LLM_NAME", "APERTUS_MODEL")
	if model == "" {
		model = ApertusModel
	}
	return base, key, model
}

// IsApertusModel accepts any Apertus identifier (8B or 70B, any casing used by
// a provider or vLLM served name) and rejects other model families.
func IsApertusModel(id string) bool {
	return modelIDPattern.MatchString(id) && strings.Contains(strings.ToLower(id), "apertus")
}

// plainHTTPAllowed limits unencrypted endpoints to loopback, single-label
// container service names and private addresses inside the operator's network.
func plainHTTPAllowed(host string) bool {
	if host == "localhost" || (host != "" && !strings.Contains(host, ".") && net.ParseIP(host) == nil) {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && (ip.IsLoopback() || ip.IsPrivate())
}

// ValidateConfiguration performs no request and never prints private configuration.
func (p *Provider) ValidateConfiguration() error {
	if p == nil || p.Key == "" {
		return errors.New("Apertus belum dikonfigurasi. Anda tetap dapat menyusun draft manual.")
	}
	if !IsApertusModel(p.Model) {
		return errors.New("Konfigurasi model harus keluarga Apertus; fallback model lain tidak diizinkan.")
	}
	u, err := url.Parse(p.BaseURL)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Scheme != "https" && !(u.Scheme == "http" && plainHTTPAllowed(u.Hostname()))) {
		return errors.New("Endpoint Apertus tidak valid.")
	}
	if p.Client == nil {
		return errors.New("Klien Apertus belum dikonfigurasi.")
	}
	return nil
}

func (p *Provider) endpoint(resource string) string {
	base := strings.TrimRight(p.BaseURL, "/")
	if !strings.HasSuffix(base, "/v1") {
		base += "/v1"
	}
	return base + "/" + resource
}

// CheckAvailability makes one read-only catalog request. A listed model is
// availability metadata, never proof that a completion used those weights.
func (p *Provider) CheckAvailability(ctx context.Context) error {
	if err := p.ValidateConfiguration(); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.endpoint("models"), nil)
	if err != nil {
		return errors.New("Tidak dapat menyiapkan pemeriksaan katalog AI.")
	}
	req.Header.Set("Authorization", "Bearer "+p.Key)
	req.Header.Set("User-Agent", "LaporSiap/0.1")
	resp, err := p.Client.Do(req)
	if err != nil {
		return errors.New("Katalog AI tidak merespons; koneksi atau batas waktu gagal.")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Katalog AI mengembalikan HTTP %d.", resp.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 256001))
	if err != nil || len(raw) > 256000 {
		return errors.New("Respons katalog AI tidak valid atau terlalu besar.")
	}
	var catalog struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if json.Unmarshal(raw, &catalog) != nil || catalog.Data == nil {
		return errors.New("Respons katalog AI tidak valid atau terlalu besar.")
	}
	for _, model := range catalog.Data {
		if strings.EqualFold(model.ID, p.Model) {
			return nil
		}
	}
	return errors.New("Model Apertus yang dipilih tidak tercantum pada katalog provider.")
}
