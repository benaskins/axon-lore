package lore_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/benaskins/axon-lore"
)

func ExampleNewWriter() {
	w := lore.NewWriter("/tmp/myproject")
	fmt.Println(w.BaseDir)
	fmt.Println(w.ADRDir)
	// Output:
	// /tmp/myproject
	// docs/adr
}

func ExampleWriter_Write() {
	dir, err := os.MkdirTemp("", "adr-example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	w := lore.NewWriter(dir)
	err = w.Write(lore.ADR{
		Number:       1,
		Title:        "Use SQLite",
		Status:       "accepted",
		Context:      "We need an embedded database.",
		Decision:     "Use SQLite for local storage.",
		Consequences: "No external database dependency.",
		Date:         "2026-04-16",
	})
	if err != nil {
		panic(err)
	}

	// Verify the file was created.
	entries, _ := os.ReadDir(filepath.Join(dir, "docs", "adr"))
	fmt.Println(entries[0].Name())
	// Output:
	// 001-use-sqlite.md
}
