package main

import (
	"context"
	"flag"
	"fmt"
	"laporsiap/internal/infrastructure/postgres"
	"os"
	"time"
)

func main() {
	purge := flag.Bool("purge", false, "Delete only expired database cases and old anonymous AI counts")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	store, err := postgres.OpenStore(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Database unavailable.")
		os.Exit(1)
	}
	defer store.Pool.Close()
	var n int64
	if *purge {
		n, err = store.PurgeExpired(ctx)
	} else {
		n, err = store.ExpiredCount(ctx)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Maintenance failed.")
		os.Exit(1)
	}
	fmt.Printf("expired_cases=%d deleted=%t\n", n, *purge)
}
