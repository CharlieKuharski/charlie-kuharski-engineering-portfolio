package protection

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type workerProtectedEntry struct {
	Mode       fs.FileMode
	Data       []byte
	LinkTarget string
}

type workerProtectedSnapshot struct {
	workspacePath         string
	allowedTests          map[string]struct{}
	mutableProtectedPaths []string
	entries               map[string]workerProtectedEntry
	protectUploads        bool
}

func snapshotWorkerProtectedPaths(workspacePath string, allowedFiles []string, protectUploads bool, mutableProtectedPaths ...string) (workerProtectedSnapshot, error) {
	workspacePath = filepath.Clean(strings.TrimSpace(workspacePath))
	if workspacePath == "" || workspacePath == "." {
		return workerProtectedSnapshot{}, fmt.Errorf("protected worker snapshot requires a workspace path")
	}
	allowedTests := map[string]struct{}{}
	for _, rel := range allowedFiles {
		rel = strings.Trim(filepath.ToSlash(filepath.Clean(strings.TrimSpace(rel))), "/")
		if strings.HasPrefix(rel, acceptanceTestWorkspaceDir+"/") {
			allowedTests[rel] = struct{}{}
		}
	}
	snapshot := workerProtectedSnapshot{
		workspacePath:         workspacePath,
		allowedTests:          allowedTests,
		mutableProtectedPaths: normalizeWorkerMutablePaths(mutableProtectedPaths),
		entries:               map[string]workerProtectedEntry{},
		protectUploads:        protectUploads,
	}
	if err := snapshot.walk(func(path string, rel string, entry fs.DirEntry) error {
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		protected := workerProtectedEntry{Mode: info.Mode()}
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			protected.LinkTarget, err = os.Readlink(path)
		case info.Mode().IsRegular():
			protected.Data, err = os.ReadFile(path)
		case info.IsDir():
		default:
			return fmt.Errorf("unsupported protected worker path type %s", rel)
		}
		if err != nil {
			return err
		}
		snapshot.entries[rel] = protected
		return nil
	}); err != nil {
		return workerProtectedSnapshot{}, fmt.Errorf("snapshot protected worker paths: %w", err)
	}
	return snapshot, nil
}

func (s workerProtectedSnapshot) Restore() ([]string, error) {
	current, err := s.currentEntries()
	if err != nil {
		return nil, err
	}
	changed := changedWorkerProtectedPaths(s.entries, current)
	if len(changed) == 0 {
		return nil, nil
	}

	removePaths := make([]string, 0, len(current))
	for rel := range current {
		if _, ok := s.entries[rel]; !ok {
			removePaths = append(removePaths, rel)
		}
	}
	slices.SortFunc(removePaths, func(a string, b string) int {
		return strings.Count(b, "/") - strings.Count(a, "/")
	})
	for _, rel := range removePaths {
		path := filepath.Join(s.workspacePath, filepath.FromSlash(rel))
		if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
			return changed, fmt.Errorf("remove unauthorized protected path %s: %w", rel, err)
		}
	}

	directories := []string{}
	files := []string{}
	for rel, entry := range s.entries {
		if entry.Mode.IsDir() {
			directories = append(directories, rel)
		} else {
			files = append(files, rel)
		}
	}
	slices.SortFunc(directories, func(a string, b string) int {
		return strings.Count(a, "/") - strings.Count(b, "/")
	})
	for _, rel := range directories {
		path := filepath.Join(s.workspacePath, filepath.FromSlash(rel))
		if info, statErr := os.Lstat(path); statErr == nil && !info.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				return changed, fmt.Errorf("replace protected directory %s: %w", rel, err)
			}
		} else if statErr != nil && !os.IsNotExist(statErr) {
			return changed, fmt.Errorf("inspect protected directory %s: %w", rel, statErr)
		}
		if err := os.MkdirAll(path, s.entries[rel].Mode.Perm()); err != nil {
			return changed, fmt.Errorf("restore protected directory %s: %w", rel, err)
		}
		if err := os.Chmod(path, s.entries[rel].Mode.Perm()); err != nil {
			return changed, fmt.Errorf("restore protected directory mode %s: %w", rel, err)
		}
	}
	slices.Sort(files)
	for _, rel := range files {
		entry := s.entries[rel]
		path := filepath.Join(s.workspacePath, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return changed, fmt.Errorf("restore protected parent %s: %w", rel, err)
		}
		if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
			return changed, fmt.Errorf("replace protected path %s: %w", rel, err)
		}
		if entry.Mode&os.ModeSymlink != 0 {
			if err := os.Symlink(entry.LinkTarget, path); err != nil {
				return changed, fmt.Errorf("restore protected symlink %s: %w", rel, err)
			}
			continue
		}
		if err := os.WriteFile(path, entry.Data, entry.Mode.Perm()); err != nil {
			return changed, fmt.Errorf("restore protected file %s: %w", rel, err)
		}
		if err := os.Chmod(path, entry.Mode.Perm()); err != nil {
			return changed, fmt.Errorf("restore protected file mode %s: %w", rel, err)
		}
	}

	verified, err := s.currentEntries()
	if err != nil {
		return changed, err
	}
	if remaining := changedWorkerProtectedPaths(s.entries, verified); len(remaining) > 0 {
		return changed, fmt.Errorf("protected worker paths did not restore: %s", strings.Join(remaining, ", "))
	}
	return changed, nil
}

func (s workerProtectedSnapshot) currentEntries() (map[string]workerProtectedEntry, error) {
	entries := map[string]workerProtectedEntry{}
	if err := s.walk(func(path string, rel string, entry fs.DirEntry) error {
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		current := workerProtectedEntry{Mode: info.Mode()}
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			current.LinkTarget, err = os.Readlink(path)
		case info.Mode().IsRegular():
			current.Data, err = os.ReadFile(path)
		case info.IsDir():
		default:
			return fmt.Errorf("unsupported protected worker path type %s", rel)
		}
		if err != nil {
			return err
		}
		entries[rel] = current
		return nil
	}); err != nil {
		return nil, fmt.Errorf("inspect protected worker paths: %w", err)
	}
	return entries, nil
}

func (s workerProtectedSnapshot) walk(visit func(path string, rel string, entry fs.DirEntry) error) error {
	return filepath.WalkDir(s.workspacePath, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == s.workspacePath {
			return nil
		}
		rel, err := filepath.Rel(s.workspacePath, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !s.protectUploads && (rel == "uploads" || strings.HasPrefix(rel, "uploads/")) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if workerResultArtifactPath(rel) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !workerPathProtected(rel, s.allowedTests, s.mutableProtectedPaths) {
			if entry.IsDir() && workerPathCanSkip(rel, s.allowedTests, s.mutableProtectedPaths) {
				return filepath.SkipDir
			}
			return nil
		}
		return visit(path, rel, entry)
	})
}

func workerPathProtected(rel string, allowedTests map[string]struct{}, mutableProtectedPaths []string) bool {
	rel = strings.Trim(filepath.ToSlash(rel), "/")
	if rel == "" {
		return false
	}
	if workerResultArtifactPath(rel) {
		return false
	}
	if _, ok := allowedTests[rel]; ok {
		return false
	}
	if workerPathAncestorOfAllowedTest(rel, allowedTests) {
		return false
	}
	if workerPathMutable(rel, mutableProtectedPaths) || workerPathAncestorOfMutable(rel, mutableProtectedPaths) {
		return false
	}
	first := strings.SplitN(rel, "/", 2)[0]
	switch first {
	case ".git", ".openmission", "tests", "uploads":
		return true
	}
	base := filepath.Base(filepath.FromSlash(rel))
	if base == ".env" || strings.HasPrefix(base, ".env.") {
		return true
	}
	switch base {
	case ".npmrc", ".yarnrc", ".netrc", "credentials.json":
		return true
	}
	switch strings.ToLower(filepath.Ext(base)) {
	case ".pem", ".key", ".p12", ".pfx":
		return true
	}
	return false
}

func workerPathCanSkip(rel string, allowedTests map[string]struct{}, mutableProtectedPaths []string) bool {
	rel = strings.Trim(filepath.ToSlash(rel), "/")
	if rel == "" {
		return false
	}
	if workerResultArtifactPath(rel) {
		return true
	}
	if workerPathAncestorOfAllowedTest(rel, allowedTests) {
		return false
	}
	if workerPathAncestorOfMutable(rel, mutableProtectedPaths) {
		return false
	}
	first := strings.SplitN(rel, "/", 2)[0]
	switch first {
	case "node_modules", "dist", "coverage", ".cache", ".vite":
		return true
	}
	return !workerPathProtected(rel, nil, nil)
}

func workerResultArtifactPath(rel string) bool {
	rel = strings.Trim(filepath.ToSlash(rel), "/")
	if rel == "" {
		return false
	}
	for _, root := range []string{
		"test-results",
		".openmission/developer-test-results",
		"playwright-report",
		"blob-report",
		"screenshots",
		"traces",
		"videos",
	} {
		if rel == root || strings.HasPrefix(rel, root+"/") {
			return true
		}
	}
	return false
}

func workerPathAncestorOfAllowedTest(rel string, allowedTests map[string]struct{}) bool {
	rel = strings.Trim(filepath.ToSlash(rel), "/")
	if rel == "" {
		return false
	}
	for allowed := range allowedTests {
		if strings.HasPrefix(allowed, rel+"/") {
			return true
		}
	}
	return false
}

func normalizeWorkerMutablePaths(paths []string) []string {
	normalized := []string{}
	for _, rel := range paths {
		rel = strings.Trim(filepath.ToSlash(filepath.Clean(strings.TrimSpace(rel))), "/")
		if rel != "" && rel != "." {
			normalized = append(normalized, rel)
		}
	}
	slices.Sort(normalized)
	return slices.Compact(normalized)
}

func workerPathMutable(rel string, mutableProtectedPaths []string) bool {
	rel = strings.Trim(filepath.ToSlash(rel), "/")
	for _, allowed := range mutableProtectedPaths {
		if rel == allowed || strings.HasPrefix(rel, allowed+"/") {
			return true
		}
	}
	return false
}

func workerPathAncestorOfMutable(rel string, mutableProtectedPaths []string) bool {
	rel = strings.Trim(filepath.ToSlash(rel), "/")
	for _, allowed := range mutableProtectedPaths {
		if strings.HasPrefix(allowed, rel+"/") {
			return true
		}
	}
	return false
}

func changedWorkerProtectedPaths(expected map[string]workerProtectedEntry, actual map[string]workerProtectedEntry) []string {
	paths := map[string]struct{}{}
	for rel := range expected {
		paths[rel] = struct{}{}
	}
	for rel := range actual {
		paths[rel] = struct{}{}
	}
	changed := []string{}
	for rel := range paths {
		left, leftOK := expected[rel]
		right, rightOK := actual[rel]
		if !leftOK || !rightOK || left.Mode != right.Mode || left.LinkTarget != right.LinkTarget || !bytes.Equal(left.Data, right.Data) {
			changed = append(changed, rel)
		}
	}
	slices.Sort(changed)
	return changed
}
