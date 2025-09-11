package blog

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func write(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestFilesStore_LoadsAndSorts(t *testing.T) {
	td := t.TempDir()
	write(t, td, "b.md", `---
title: "B"
date: 2025-08-02
---
b`)
	write(t, td, "a.md", `---
title: "A"
date: 2025-08-01
---
a`)

	s, err := NewFilesStore(td)
	if err != nil {
		t.Fatal(err)
	}
	got := s.All()
	if len(got) != 2 {
		t.Fatalf("posts=%d, want 2", len(got))
	}
	if got[0].Title != "B" || got[1].Title != "A" {
		t.Fatalf("order wrong: %+v", []string{got[0].Title, got[1].Title})
	}
	if _, ok := s.BySlug(Slugify("B")); !ok {
		t.Fatalf("BySlug(B) not found")
	}
}

func TestFilesStore_DraftAndFuture(t *testing.T) {
	td := t.TempDir()
	write(t, td, "draft.md", `---
title: "Draft"
date: 2025-08-01
draft: true
---
x`)
	future := time.Now().Add(48 * time.Hour).Format("2006-01-02")
	write(t, td, "future.md", `---
title: "Future"
date: `+future+`
---
x`)

	s, err := NewFilesStore(td) // showDrafts=false by default
	if err != nil {
		t.Fatal(err)
	}
	got := s.All()
	if len(got) != 0 {
		t.Fatalf("posts=%d, want 0 (draft+future filtered)", len(got))
	}

	// WithDrafts=true keeps draft and future
	s2, err := NewFilesStore(td, WithDrafts(true))
	if err != nil {
		t.Fatal(err)
	}
	got2 := s2.All()
	if len(got2) != 2 {
		t.Fatalf("posts=%d, want 2 with drafts", len(got2))
	}
}

func TestFilesStore_DuplicateSlugs(t *testing.T) {
	td := t.TempDir()
	write(t, td, "one.md", `---
title: "Same Title"
date: 2025-08-01
---
x`)
	write(t, td, "two.md", `---
title: "Same Title"
date: 2025-08-02
---
y`)

	s, err := NewFilesStore(td)
	if err != nil {
		t.Fatal(err)
	}
	all := s.All()
	if all[0].Slug == all[1].Slug {
		t.Fatalf("slugs not unique: %q", all[0].Slug)
	}
}

