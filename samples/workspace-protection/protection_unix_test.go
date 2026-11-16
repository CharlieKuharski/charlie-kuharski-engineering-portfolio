//go:build unix

package protection

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestWorkerWorkspaceProtectionRestoresRegularFileContentAndModeUnderUmask(t *testing.T) {
	workspacePath := receiptTestWorkspace(t, "protected-file-mode")
	storyPath := filepath.Join(workspacePath, ".openmission", "story.json")
	resultPath := filepath.Join(workspacePath, "test-results", "ac-1", "result.json")
	writeFileForTest(t, storyPath, "approved story\n")
	writeFileForTest(t, resultPath, "original result\n")
	if err := os.Chmod(storyPath, 0o664); err != nil {
		t.Fatal(err)
	}

	snapshot, err := snapshotWorkerProtectedPaths(workspacePath, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	writeFileForTest(t, storyPath, "worker story\n")
	if err := os.Chmod(storyPath, 0o600); err != nil {
		t.Fatal(err)
	}
	writeFileForTest(t, resultPath, "worker result\n")

	previousUmask := syscall.Umask(0o022)
	defer syscall.Umask(previousUmask)
	changed, err := snapshot.Restore()
	if err != nil {
		t.Fatal(err)
	}
	if len(changed) == 0 {
		t.Fatal("protected Story mutation was not detected")
	}
	assertCheckpointFixtureFile(t, workspacePath, ".openmission/story.json", "approved story\n")
	assertCheckpointFixtureFile(t, workspacePath, "test-results/ac-1/result.json", "worker result\n")
	info, err := os.Stat(storyPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o664 {
		t.Fatalf("restored Story mode = %04o, want 0664", got)
	}
}
