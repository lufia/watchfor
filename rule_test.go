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
		Ok     bool
		Expect string
	}{
		{".c", "test.c", true, "test.rc"},         // 返還前のほうが短い
		{".go", "test.go", true, "test.rc"},       // 返還前と同じ長さ
		{".cpp", "test.cpp", true, "test.rc"},     // 返還前のほうが長い
		{".", "test.", true, "test.rc"},           // .で終わる
		{"", "file", true, "file.rc"},             // 拡張子なし
		{".c", "dir/test.c", true, "dir/test.rc"}, // ディレクトリ有り
		{".c", "test.g", false, "test.c"},         // マッチしない
	}
	for _, r := range tab {
		rule := &Rule{
			SrcExt:  FileExt(r.SrcExt),
			DestExt: destExt,
		}
		s, ok := rule.Convert(r.Path)
		if ok != r.Ok {
			t.Errorf("Convert(%v) = _, %v; want _, %v", r.Path, ok, r.Ok)
		}
		if r.Ok && s != r.Expect {
			t.Errorf("Convert(%v) = %v, _; want %v, _", r.Path, s, r.Expect)
		}
	}
}
