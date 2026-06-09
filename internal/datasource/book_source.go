package datasource

import "context"

// BookSummary 书籍摘要（通用）
type BookSummary struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Rating   string `json:"rating"`
	Year     string `json:"year"`
	CoverURL string `json:"cover_url,omitempty"`
}

// BookDetail 书籍详情（通用）
type BookDetail struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Rating    string `json:"rating"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	ISBN      string `json:"isbn"`
	Intro     string `json:"intro"`
	URL       string `json:"url"`
	CoverURL  string `json:"cover_url"`
}

// BookSource 通用书籍数据源接口
type BookSource interface {
	Name() string
	Search(ctx context.Context, keyword string, limit int) ([]BookSummary, error)
	GetBookDetail(ctx context.Context, id string) (*BookDetail, error)
	DownloadCover(ctx context.Context, id, coverURL string) error
}
