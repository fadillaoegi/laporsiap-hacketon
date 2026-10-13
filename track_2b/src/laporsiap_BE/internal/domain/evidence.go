package domain

import "encoding/json"

type AttachmentData struct {
	Meta  Attachment
	Bytes []byte
}

type HistoryEntry struct {
	Revision  int64           `json:"revision"`
	Action    string          `json:"action"`
	CreatedAt string          `json:"created_at"`
	Document  json.RawMessage `json:"document"`
}
