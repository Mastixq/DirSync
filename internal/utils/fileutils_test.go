package fileutils

import (
	"bytes"
	"path/filepath"
	"testing"
	"time"
)

func TestDeleteMissing_RemovesExtraFiles(t *testing.T) {
	fs := NewFakeFS()

	absSrc, err := filepath.Abs("src")
	if err != nil {
		t.Fatalf("filepath.Abs failed: %v", err)
	}
	absDst, err := filepath.Abs("dst")
	if err != nil {
		t.Fatalf("filepath.Abs failed: %v", err)
	}

	fs.MkdirAll(absSrc, 0755)
	fs.MkdirAll(absDst, 0755)

	fs.files[filepath.Join(absSrc, "file.txt")] = &fakeFile{
		name:    "file.txt",
		data:    bytes.NewBufferString("hello"),
		mode:    0644,
		modTime: time.Now(),
	}

	fs.files[filepath.Join(absDst, "file.txt")] = &fakeFile{
		name:    "file.txt",
		data:    bytes.NewBufferString("obsolete"),
		mode:    0644,
		modTime: time.Now(),
	}
	fs.files[filepath.Join(absDst, "extra.txt")] = &fakeFile{
		name:    "extra.txt",
		data:    bytes.NewBufferString("to delete"),
		mode:    0644,
		modTime: time.Now(),
	}

	err = DeleteMissing(fs, absSrc, absDst)
	if err != nil {
		t.Fatalf("DeleteMissing failed: %v", err)
	}

	if _, exists := fs.files[filepath.Join(absDst, "extra.txt")]; exists {
		t.Errorf("File dst/extra.txt was not removed")
	}
}

func TestFilesAreEqual(t *testing.T) {
	type testCase struct {
		name     string
		contentA string
		contentB string
		modTimeA time.Time
		modTimeB time.Time
		want     bool
	}

	cases := []testCase{
		{
			name:     "identical content and modtime",
			contentA: "abc",
			contentB: "abc",
			modTimeA: time.Now(),
			modTimeB: time.Now(),
			want:     true,
		},
		{
			name:     "different content",
			contentA: "abc",
			contentB: "xyz",
			modTimeA: time.Now(),
			modTimeB: time.Now(),
			want:     false,
		},
		{
			name:     "same content, different modtime",
			contentA: "same",
			contentB: "same",
			modTimeA: time.Now(),
			modTimeB: time.Now().Add(time.Second),
			want:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fs := NewFakeFS()

			fs.files["a.txt"] = &fakeFile{
				name:    "a.txt",
				data:    bytes.NewBufferString(tc.contentA),
				mode:    0644,
				modTime: tc.modTimeA,
			}
			fs.files["b.txt"] = &fakeFile{
				name:    "b.txt",
				data:    bytes.NewBufferString(tc.contentB),
				mode:    0644,
				modTime: tc.modTimeB,
			}

			got, err := filesAreEqual(fs, "a.txt", "b.txt")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestCopyFilePreserveTree_CopiesFileCorrectly(t *testing.T) {
	fs := NewFakeFS()
	src := "src/file.txt"
	dst := "dst"
	content := "hello world"

	fs.MkdirAll("src", 0755)
	fs.MkdirAll("dst", 0755)
	fs.files[src] = &fakeFile{
		name:    "file.txt",
		data:    bytes.NewBufferString(content),
		mode:    0644,
		modTime: time.Now(),
	}

	err := CopyFilePreserveTree(fs, src, "src", dst)
	if err != nil {
		t.Fatalf("CopyFilePreserveTree failed: %v", err)
	}

	targetPath := dst + "/file.txt"
	copied, exists := fs.files[targetPath]
	if !exists {
		t.Fatalf("Expected file %s to exist", targetPath)
	}

	if copied.data.String() != content {
		t.Errorf("File content mismatch. Got: %s, Want: %s", copied.data.String(), content)
	}
}
