package protection

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkerWorkspaceProtectionSkipsResultArtifacts(t *testing.T) {
	workspacePath := receiptTestWorkspace(t, "result-artifact-snapshot")
	writeFileForTest(t, filepath.Join(workspacePath, ".openmission/story.json"), "story\n")

	resultFiles := []string{
		"test-results/ac-1/controller/playwright-report/index.html",
		".openmission/developer-test-results/ac-1/latest/report/index.html",
		"playwright-report/index.html",
		"blob-report/report.zip",
		"screenshots/failure.png",
		"traces/trace.zip",
		"videos/run.webm",
	}
	largeArtifact := strings.Repeat("diagnostic-output\n", 32*1024)
	for _, rel := range resultFiles {
		writeFileForTest(t, filepath.Join(workspacePath, filepath.FromSlash(rel)), largeArtifact)
	}

	snapshot, err := snapshotWorkerProtectedPaths(workspacePath, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	for rel := range snapshot.entries {
		if workerResultArtifactPath(rel) {
			t.Fatalf("result artifact was snapshotted: %s", rel)
		}
	}

	for _, rel := range resultFiles[2:] {
		writeFileForTest(t, filepath.Join(workspacePath, filepath.FromSlash(rel)), "regenerated\n")
	}
	deletedPath := filepath.Join(workspacePath, filepath.FromSlash(resultFiles[0]))
	if err := os.Remove(deletedPath); err != nil {
		t.Fatal(err)
	}
	replacedPath := filepath.Join(workspacePath, filepath.FromSlash(resultFiles[1]))
	if err := os.Remove(replacedPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(replacedPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFileForTest(t, filepath.Join(workspacePath, "test-results/ac-2/new-result.json"), "new result\n")
	changed, err := snapshot.Restore()
	if err != nil {
		t.Fatal(err)
	}
	if len(changed) != 0 {
		t.Fatalf("result-only churn was treated as a protected change: %v", changed)
	}
	if _, err := os.Stat(deletedPath); !os.IsNotExist(err) {
		t.Fatalf("deleted result artifact was restored: %v", err)
	}
	if info, err := os.Stat(replacedPath); err != nil || !info.IsDir() {
		t.Fatalf("replaced result artifact was changed: info=%v err=%v", info, err)
	}
	assertCheckpointFixtureFile(t, workspacePath, "test-results/ac-2/new-result.json", "new result\n")
	for _, rel := range resultFiles[2:] {
		assertCheckpointFixtureFile(t, workspacePath, rel, "regenerated\n")
	}
}

func TestWorkerWorkspaceProtectionFailsClosedWhenContractCannotRestore(t *testing.T) {
	workspacePath := receiptTestWorkspace(t, "unrestorable-contract")
	writeFileForTest(t, filepath.Join(workspacePath, ".openmission/story.json"), "approved story\n")
	snapshot, err := snapshotWorkerProtectedPaths(workspacePath, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(workspacePath); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(workspacePath, []byte("not a workspace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := snapshot.Restore(); err == nil {
		t.Fatal("expected protected-contract restoration failure")
	}
}
