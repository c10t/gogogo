package backup

import (
	"os"
	"testing"
)

func setup(t *testing.T) {
	t.Log("--- creating directory for test ---")
	os.MkdirAll("test/output", 0777)
}
func teardown(t *testing.T) {
	t.Log("--- removing directory for test ---")
	os.RemoveAll("test/output")
}

func TestZipArchive(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := ZIP.Archive("test/hash1", "test/output/1.zip")
	if err != nil {
		t.Fatalf("ZIP.Archive() failed")
	}

	t.Logf("TODO: implement Unarchive")
}
