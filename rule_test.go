package main

import (
	"testing"
)

func TestFileExtReplace(t *testing.T) {
	ext := FileExt(".rc")
	tab := []struct {
		Path   string
		Expect string
	}{
		{"test.c", "test.rc"},         // 返還前のほうが短い
		{"test.go", "test.rc"},        // 返還前と同じ長さ
		{"test.cpp", "test.rc"},       // 返還前のほうが長い
		{"test.", "test.rc"},          // .で終わる
		{"file", "file"},              // 拡張子なし
		{"dir/test.c", "dir/test.rc"}, // ディレクトリ有り
	}
	for _, r := range tab {
		s := ext.Replace(r.Path)
		if s != r.Expect {
			t.Errorf("Replace(%v) = %v; want %v", r.Path, s, r.Expect)
		}
	}
}
