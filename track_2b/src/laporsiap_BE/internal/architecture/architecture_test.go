package architecture

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestDependenciesPointInward(t *testing.T) {
	root := ".."
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		layer := strings.Split(filepath.ToSlash(rel), "/")[0]
		f, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, i := range f.Imports {
			dep, _ := strconv.Unquote(i.Path.Value)
			internal := strings.TrimPrefix(dep, "laporsiap/internal/")
			isInternal := internal != dep
			external := strings.Contains(strings.Split(dep, "/")[0], ".")
			allowed := true
			switch layer {
			case "domain":
				allowed = !isInternal && !external && dep != "net/http" && !strings.Contains(dep, "pgx") && !strings.Contains(dep, "fpdf")
			case "application":
				allowed = (!isInternal && !external && dep != "net/http" && !strings.Contains(dep, "pgx") && !strings.Contains(dep, "fpdf")) || (isInternal && (internal == "domain" || strings.HasPrefix(internal, "application")))
			case "transport":
				allowed = (!isInternal && !external) || internal == "domain" || internal == "application" || strings.HasPrefix(internal, "transport/")
			case "infrastructure":
				allowed = !isInternal || internal == "domain" || internal == "application" || strings.HasPrefix(internal, "infrastructure/")
			}
			if !allowed {
				t.Errorf("%s imports outward dependency %s", rel, dep)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
