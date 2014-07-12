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

func TestRuleReplace(t *testing.T) {
	destExt := FileExt(".rc")
	tab := []struct {
		SrcExt string
		Path   string
		Expect string
	}{
		{".c", "test.c", "test.rc"},         // 返還前のほうが短い
		{".go", "test.go", "test.rc"},       // 返還前と同じ長さ
		{".cpp", "test.cpp", "test.rc"},     // 返還前のほうが長い
		{".", "test.", "test.rc"},           // .で終わる
		{"", "file", "file.rc"},             // 拡張子なし
		{".c", "dir/test.c", "dir/test.rc"}, // ディレクトリ有り
		{".c", "test.g", "test.g"},          // マッチしない
	}
	for _, r := range tab {
		rule := &Rule{
			SrcExt:  FileExt(r.SrcExt),
			DestExt: destExt,
		}
		s := rule.ConvertFilename(r.Path)
		if s != r.Expect {
			t.Errorf("Convert(%v) = %v; want %v", r.Path, s, r.Expect)
		}
	}
}
