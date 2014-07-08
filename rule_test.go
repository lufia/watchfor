package main

import (
	"testing"
)

func TestFileExtMatch(t *testing.T) {
	ext := FileExt(".c")
	tab := []struct {
		Path   string
		Expect bool
	}{
		{"test.c", true},      // マッチ
		{".c", true},          // マッチ
		{"test.cc", false},    // 文字が多い
		{"test.c.x", false},   // ファイル名の最後ではない
		{"a.c/test.x", false}, // 最後のパスではない
		{"test.", false},      // 拡張子が無い
	}
	for _, r := range tab {
		if v := ext.Match(r.Path); v != r.Expect {
			t.Errorf("Match(%v) = %v; want %v", r.Path, v, r.Expect)
		}
	}
}

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
