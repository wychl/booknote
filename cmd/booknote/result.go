// result.go
package main

import (
	"encoding/json"
	"os"

	"github.com/wychl/booknote/internal/datasource"
)

// Result 命令执行后输出到 stdout 的 JSON 结构
type Result struct {
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
	HTML     string `json:"html,omitempty"`
	Image    string `json:"image,omitempty"`
	JSONFile string `json:"json_file,omitempty"`
	BookNote
}

func (r *Result) output() {
	_ = json.NewEncoder(os.Stdout).Encode(r)
}

// BookNote 笔记数据结构
type BookNote struct {
	Book  datasource.BookDetail `json:"book"`
	Note  string                `json:"note"`
	Theme string                `json:"theme"`
	Style string                `json:"style"`
	Tags  []string              `json:"tags,omitempty"`
	Bgm   []string              `json:"bgm,omitempty"`
}
