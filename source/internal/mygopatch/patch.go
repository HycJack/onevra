// Package mygopatch applies the narrowly scoped, pinned MyGo integration patch.
// It never modifies the module cache or downloads an unpinned source tree.
package mygopatch

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const Module = "github.com/egoist/mygo"
const Version = "v0.2.15"

//go:embed changes.patch
var Patch []byte

// Apply writes only changed files into dst, reading original files from base.
// The returned map identifies changed paths. Context mismatches fail closed;
// callers must pass a writable destination, never the downloaded module cache.
func Apply(base, dst string, patch []byte) (map[string]string, error) {
	s := bufio.NewScanner(bytes.NewReader(patch))
	s.Buffer(make([]byte, 4096), 4<<20)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	files := map[string]string{}
	for i := 0; i < len(lines); {
		if !strings.HasPrefix(lines[i], "--- ") {
			return nil, fmt.Errorf("invalid patch header at line %d", i+1)
		}
		isNew := lines[i] == "--- /dev/null"
		i++
		if i >= len(lines) || !strings.HasPrefix(lines[i], "+++ b/") {
			return nil, fmt.Errorf("invalid new path")
		}
		rel := strings.TrimPrefix(lines[i], "+++ b/")
		i++
		if !fs.ValidPath(rel) || strings.Contains(rel, ":") || strings.Contains(rel, "\\") {
			return nil, fmt.Errorf("unsafe patch path %q", rel)
		}
		path := filepath.Join(base, filepath.FromSlash(rel))
		var original []byte
		if !isNew {
			var err error
			original, err = os.ReadFile(path)
			if err != nil {
				return nil, err
			}
		}
		old := strings.Split(strings.TrimSuffix(string(original), "\n"), "\n")
		if len(original) == 0 {
			old = nil
		}
		var out []string
		pos := 0
		for i < len(lines) && strings.HasPrefix(lines[i], "@@ ") {
			fields := strings.Fields(lines[i])
			i++
			if len(fields) < 4 || fields[3] != "@@" {
				return nil, fmt.Errorf("invalid hunk in %s", rel)
			}
			start, count, err := span(fields[1], '-')
			if err != nil {
				return nil, err
			}
			_, newCount, err := span(fields[2], '+')
			if err != nil {
				return nil, err
			}
			at := max(start-1, 0)
			if at < pos || at > len(old) {
				return nil, fmt.Errorf("hunk order mismatch in %s", rel)
			}
			out = append(out, old[pos:at]...)
			pos = at
			removed, added := 0, 0
			for removed < count || added < newCount {
				if i >= len(lines) || len(lines[i]) == 0 {
					return nil, fmt.Errorf("truncated hunk in %s", rel)
				}
				line := lines[i]
				i++
				if line[0] != ' ' && line[0] != '-' && line[0] != '+' {
					return nil, fmt.Errorf("invalid hunk in %s", rel)
				}
				if line[0] != '+' {
					if pos >= len(old) || old[pos] != line[1:] {
						return nil, fmt.Errorf("MyGo patch context mismatch: %s:%d", rel, pos+1)
					}
					pos++
					removed++
				}
				if line[0] != '-' {
					out = append(out, line[1:])
					added++
				}
				if removed > count || added > newCount {
					return nil, fmt.Errorf("hunk size mismatch in %s", rel)
				}
			}
		}
		out = append(out, old[pos:]...)
		if _, duplicate := files[path]; duplicate {
			return nil, fmt.Errorf("duplicate patch path %s", rel)
		}
		target := filepath.Join(dst, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
			return nil, err
		}
		if err := os.WriteFile(target, []byte(strings.Join(out, "\n")+"\n"), 0600); err != nil {
			return nil, err
		}
		files[path] = target
	}
	return files, nil
}
func span(s string, prefix byte) (int, int, error) {
	if len(s) < 2 || s[0] != prefix {
		return 0, 0, fmt.Errorf("invalid hunk span %q", s)
	}
	parts := strings.Split(s[1:], ",")
	start, err := strconv.Atoi(parts[0])
	if err != nil || start < 0 {
		return 0, 0, fmt.Errorf("invalid hunk span %q", s)
	}
	count := 1
	if len(parts) == 2 {
		count, err = strconv.Atoi(parts[1])
	}
	if err != nil || count < 0 || len(parts) > 2 {
		return 0, 0, fmt.Errorf("invalid hunk span %q", s)
	}
	return start, count, nil
}
