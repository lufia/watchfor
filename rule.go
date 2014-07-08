package main

import (
	"os"
	"path/filepath"
	"strings"
)

// ファイル変更時に実行されるルール
type Rule struct {
	SrcExt  FileExt
	DestExt FileExt
	Cmd     Command
}

// ルール定義上のコマンド。
type Command string

// pathの拡張子をターゲットの形に変換する。
// 拡張子がなければ何もしない。
func (rule *Rule) Convert(path string) (target string, ok bool) {
	if !rule.SrcExt.Match(path) {
		return path, false
	}
	i := strings.LastIndex(path, rule.SrcExt.String())
	if i < 0 {
		panic("no file extension")
	}
	return path[0:i] + rule.DestExt.String(), true
}

// コマンドを実行してファイル生成
func (rule *Rule) Exec(src, dest string) error {
	if err := os.Setenv("source", src); err != nil {
		return err
	}
	if err := os.Setenv("target", dest); err != nil {
		return err
	}
	return System(rule.Cmd)
}

// 拡張子(.を含む)
type FileExt string

func (m FileExt) Match(file string) bool {
	return filepath.Ext(file) == m.String()
}

func (m FileExt) String() string {
	if len(m) > 0 && m[0] != '.' {
		return "." + string(m)
	} else {
		return string(m)
	}
}
