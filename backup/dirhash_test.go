package backup_test

import (
	"testing"

	"github.com/c10t/gogogo/backup"
)

func TestDirHash(t *testing.T) {
	hash1a, err := backup.DirHash("test/hash1")
	if err != nil {
		t.Fatalf("DirHash should return no error")
	}

	hash1b, err := backup.DirHash("test/hash1")
	if err != nil {
		t.Fatalf("DirHash should return no error")
	}

	if hash1a != hash1b {
		t.Fatalf("hash1a and hash1b should be identical")
	}

	hash2, err := backup.DirHash("test/hash2")
	if err != nil {
		t.Fatalf("DirHash should return no error")
	}

	if hash1a == hash2 {
		t.Fatalf("hash1a and hash2 should not be the same")
	}
}
