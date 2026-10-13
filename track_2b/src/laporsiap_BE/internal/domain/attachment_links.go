package domain

import (
	"errors"
)

// Links are user assertions, not OCR results or proof of authenticity.
func ValidateLinks(c Case, sourceIDs, transferIDs []string) error {
	if len(sourceIDs) > len(c.Messages) || len(transferIDs) > len(c.Transfers) {
		return errors.New("Tautan bukti harus menunjuk sumber dan transfer dalam kasus ini, tanpa duplikasi.")
	}
	sources, transfers := map[string]bool{}, map[string]bool{}
	for _, m := range c.Messages {
		sources[m.ID] = true
	}
	for _, t := range c.Transfers {
		transfers[t.ID] = true
	}
	for _, group := range []struct {
		ids     []string
		allowed map[string]bool
	}{{sourceIDs, sources}, {transferIDs, transfers}} {
		seen := map[string]bool{}
		for _, id := range group.ids {
			if !group.allowed[id] || seen[id] {
				return errors.New("Tautan bukti harus menunjuk sumber dan transfer dalam kasus ini, tanpa duplikasi.")
			}
			seen[id] = true
		}
	}
	return nil
}

func PruneTransferLinks(c *Case) {
	valid := map[string]bool{}
	for _, t := range c.Transfers {
		valid[t.ID] = true
	}
	for i := range c.Attachments {
		ids := []string{}
		for _, id := range c.Attachments[i].TransferIDs {
			if valid[id] {
				ids = append(ids, id)
			}
		}
		c.Attachments[i].TransferIDs = ids
	}
}
