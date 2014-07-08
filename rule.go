package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ファイル変更時に実行されるルール
type Rule struct {
	SrcExt  FileExt
	DestExt FileExt
	Command Command
}

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

// ルール定義上のコマンド。
// $varなど、シェル変数もそのまま保持する。
type Command []string

func (cmd Command) ExpandSrc(src SourceFile) []string {
	args := make([]string, len(cmd))
	for i, arg := range cmd {
		args[i] = src.Expand(arg)
	}
	return args
}

// 変更があったソースファイル
type SourceFile string

// $var, ${var}等を展開した文字列を返す
func (src SourceFile) Expand(s string) string {
	return os.Expand(s, src.getenv)
}

func (src SourceFile) getenv(key string) string {
	switch key {
	case "source":
		return string(src)
	default:
		return os.Getenv(key)
	}
}

func (rule *Rule) Exec(file string) error {
	if !rule.SrcExt.Match(file) {
		return nil
	}

	src := SourceFile(file)
	args := rule.Command.ExpandSrc(src)
	cmd := exec.Command(args[0], args[1:]...)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
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
