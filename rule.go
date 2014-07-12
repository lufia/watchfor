package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ルールにマッチしないファイルだった場合
	ErrNotCovered = errors.New("not covered")
)

// ファイル変更時に実行されるルール
type Rule struct {
	SrcExt  FileExt // 監視するファイルの拡張子
	DestExt FileExt // コマンド実行後のファイル拡張子
	Cmd     Command // SrcExtに変更があった場合実行するコマンド
}

func (rule *Rule) Eval(path string) (target string, err error) {
	if !rule.SrcExt.Match(path) {
		err = ErrNotCovered
		return
	}
	target = rule.ConvertFilename(path)
	if err = rule.execute(path, target); err != nil {
		return
	}
	return
}

// ルール定義上のコマンド。
type Command string

// pathの拡張子をターゲットの形に変換する。
// 拡張子がなければ何もしない。
func (rule *Rule) ConvertFilename(path string) string {
	if !rule.SrcExt.Match(path) {
		return path
	}
	i := strings.LastIndex(path, rule.SrcExt.String())
	if i < 0 {
		panic("no file extension")
	}
	return path[0:i] + rule.DestExt.String()
}

// コマンドを実行してファイル生成
func (rule *Rule) execute(src, dest string) error {
	if rule.Cmd == "" {
		return nil
	}
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
