package llm

// NoteContent 表示 AI 生成的读书笔记内容
type NoteContent struct {
	// MainText 笔记正文，建议 280~320 字
	MainText string `json:"main_text"`
}
