// Synthetic-only integration probe. All API/AI calls use Go.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	live := flag.Bool("live", false, "Send three synthetic source messages to configured Apertus via local backend")
	flag.Parse()
	if !*live {
		fmt.Println("Use -live to perform one synthetic Apertus integration call. No research batch is repeated.")
		return
	}
	client := &http.Client{Timeout: 55 * time.Second}
	var token string
	send := func(path string, body any) []byte {
		data, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "http://127.0.0.1:8787/api"+path, bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Local API request failed.")
			os.Exit(1)
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(io.LimitReader(resp.Body, 256000))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Response read failed.")
			os.Exit(1)
		}
		if resp.StatusCode >= 400 {
			var msg struct {
				Error string `json:"error"`
			}
			_ = json.Unmarshal(b, &msg)
			fmt.Fprintf(os.Stderr, "Integration probe failed (HTTP %d): %s\n", resp.StatusCode, msg.Error)
			os.Exit(1)
		}
		return b
	}
	b := send("/cases", map[string]any{"title": "Uji integrasi Apertus — sintetis", "synthetic": true, "parser": "plain", "raw_chat": "Transfer Rp275.000 saya gagal. Saldo tidak berkurang.\nCoba kirim lagi ke DEMO-REK-H03.\nSaya tidak mencoba lagi dan belum ada transaksi berhasil.", "transfers": []any{}})
	var created struct {
		Token string `json:"token"`
		Case  struct {
			ID       string `json:"id"`
			Revision int64  `json:"revision"`
		} `json:"case"`
	}
	if json.Unmarshal(b, &created) != nil || created.Token == "" {
		fmt.Fprintln(os.Stderr, "Invalid create response.")
		os.Exit(1)
	}
	token = created.Token
	b = send("/cases/"+created.Case.ID+"/draft", map[string]any{"revision": created.Case.Revision, "mode": "apertus", "external_consent": true})
	var result struct {
		Case struct {
			Revision int64           `json:"revision"`
			Draft    json.RawMessage `json:"draft"`
		} `json:"case"`
	}
	if json.Unmarshal(b, &result) != nil {
		fmt.Fprintln(os.Stderr, "Invalid draft response.")
		os.Exit(1)
	}
	safe, _ := json.MarshalIndent(map[string]any{"synthetic": true, "revision": result.Case.Revision, "draft": result.Case.Draft, "note": "One integration probe only; no human review confirmation or export performed."}, "", "  ")
	fmt.Println(string(safe))
}
