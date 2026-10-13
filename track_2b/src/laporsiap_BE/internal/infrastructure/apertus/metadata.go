package apertus

import (
	"laporsiap/internal/application"
	"net/url"
)

func (p *Provider) Info() application.ProviderInfo {
	label := "Configured Apertus provider"
	if p != nil {
		u, err := url.Parse(p.BaseURL)
		if err == nil {
			if u.Scheme == "https" && u.Hostname() == "api.publicai.co" {
				label = "Public AI"
			}
			if u.Scheme == "http" && plainHTTPAllowed(u.Hostname()) {
				label = "Local Apertus server"
			}
		}
	}
	return application.ProviderInfo{Label: label, PromptVersion: PromptVersion, PromptSHA256: PromptSHA256()}
}

var _ application.Generator = (*Provider)(nil)
