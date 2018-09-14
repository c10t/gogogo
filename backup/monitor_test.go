package backup

import (
	"fmt"
	"testing"
)

type call struct {
	Src  string
	Dest string
}

type TestArchiver struct {
	Archives []*call
	Restores []*call
}

func (a *TestArchiver) DestFmt() func(int64) string {
	return func(i int64) string {
		return fmt.Sprintf("%d.zip", i)
	}
}

func (a *TestArchiver) Archive(src, dest string) error {
	a.Archives = append(a.Archives, &call{Src: src, Dest: dest})
	return nil
}
func (a *TestArchiver) Restore(src, dest string) error {
	a.Restores = append(a.Restores, &call{Src: src, Dest: dest})
	return nil
}

func TestMonitor(t *testing.T) {
	a := &TestArchiver{}
	m := &Monitor{
		Destination: "test/archive",
		Paths: map[string]string{
			"test/hash1": "abc",
			"test/hash2": "def",
		},
		Archiver: a,
	}

	n, err := m.Now()
	if err != nil {
		t.Fatalf("monitor.Now() should return no error")
	}
	if n != 2 {
		t.Fatalf("monitor.Now() returns wrong value, expect: 2, actual: %v", n)
	}
	if len(a.Archives) != 2 {
		t.Fatalf("monitor.Archiver.Archives returns wrong value, expect: 2, actual: %v", n)
	}

	/* this test may depend on host os
	for _, call := range a.Archives {
		if !strings.HasPrefix(call.Dest, m.Destination) {
			t.Fatalf("Dest should have prefix %v but actual %v", call.Dest, m.Destination)
		}
		if !strings.HasSuffix(call.Dest, ".zip") {
			t.Fatalf("Dest should have suffix .zip but actual %v", call.Dest)
		}
	}
	*/
}
