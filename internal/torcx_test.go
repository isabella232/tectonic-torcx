package internal

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestTorcxGC(t *testing.T) {
	storeDir, err := ioutil.TempDir("", ".torcx-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(storeDir)

	dirs := []string{"v0", "v1", "v2", "v3"}
	for _, d := range dirs {
		if err := os.Mkdir(filepath.Join(storeDir, d), 0755); err != nil {
			t.Fatal(err)
		}
		touch(t, filepath.Join(storeDir, d, "a"))
	}
	touch(t, filepath.Join(storeDir, "a"))

	a, err := NewApp(Config{
		torcxStoreDir: storeDir,
		TorcxBin:      "/bin/true",
	})
	if err != nil {
		t.Fatal(err)
	}
	a.OSVersions = []string{"v2", "v3"}

	err = a.TorcxGC()
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"a", "v2/", "v3/"}
	actual := listDir(t, storeDir)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected store dir of %q, got %q", expected, actual)
	}
}

func touch(t *testing.T, path string) {
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
}

func listDir(t *testing.T, path string) []string {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() {
			n += "/"
		}
		out = append(out, n)
	}

	return out
}
