package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type bookDoubanSource struct {
	httpClient *http.Client
	cookie     string
}

// NewBookDoubanSource 创建豆瓣数据源实例
func NewBookDoubanSource(cookie string) BookSource {
	return &bookDoubanSource{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cookie:     cookie,
	}
}

func (d *bookDoubanSource) Name() string { return "douban" }

// Search 搜索书籍
func (d *bookDoubanSource) Search(ctx context.Context, keyword string, limit int) ([]BookSummary, error) {
	if limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	baseURL := "https://m.douban.com/rexxar/api/v2/search"
	params := url.Values{}
	params.Set("q", keyword)
	params.Set("type", "")
	params.Set("loc_id", "")
	params.Set("start", "0")
	params.Set("count", fmt.Sprintf("%d", limit))
	params.Set("sort", "relevance")
	params.Set("ck", "y4pS")
	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	d.setSearchHeaders(req, keyword)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResp struct {
		Subjects *struct {
			Items []struct {
				Target struct {
					ID     string `json:"id"`
					Title  string `json:"title"`
					Rating struct {
						Value float64 `json:"value"`
					} `json:"rating"`
					CoverURL string `json:"cover_url"`
					Year     string `json:"year"`
				} `json:"target"`
			} `json:"items"`
		} `json:"subjects"`
	}
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, err
	}
	if searchResp.Subjects == nil || len(searchResp.Subjects.Items) == 0 {
		return nil, fmt.Errorf("未找到相关书籍: %s", keyword)
	}

	summaries := make([]BookSummary, 0, len(searchResp.Subjects.Items))
	for _, item := range searchResp.Subjects.Items {
		summaries = append(summaries, BookSummary{
			ID:       item.Target.ID,
			Title:    item.Target.Title,
			Rating:   fmt.Sprintf("%.1f", item.Target.Rating.Value),
			Year:     item.Target.Year,
			CoverURL: item.Target.CoverURL,
		})
	}
	return summaries, nil
}

func (d *bookDoubanSource) setSearchHeaders(req *http.Request, keyword string) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://www.douban.com/search?source=suggest&q="+url.QueryEscape(keyword))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
	if d.cookie != "" {
		req.Header.Set("Cookie", d.cookie)
	}
}

// GetBookDetail 获取书籍详细信息
func (d *bookDoubanSource) GetBookDetail(ctx context.Context, id string) (*BookDetail, error) {
	urlStr := fmt.Sprintf("https://book.douban.com/subject/%s/", id)
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	d.setDetailHeaders(req)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	info := &BookDetail{ID: id, URL: urlStr}
	info.Title = strings.TrimSpace(doc.Find(`h1 span[property="v:itemreviewed"]`).Text())
	info.Rating = strings.TrimSpace(doc.Find(`.rating_num`).Text())
	info.Author = strings.TrimSpace(doc.Find(`.pl:contains("作者")`).Parent().Find("a").First().Text())
	info.Publisher = strings.TrimSpace(doc.Find(`#info .pl:contains("出版社:")`).Next().Text())
	if metaIsbn := doc.Find(`meta[property="book:isbn"]`); metaIsbn.Length() > 0 {
		info.ISBN, _ = metaIsbn.Attr("content")
	}
	if introFull := doc.Find(`#link-report .all .intro`).First(); introFull.Length() > 0 {
		info.Intro = strings.TrimSpace(introFull.Text())
	} else if introShort := doc.Find(`#link-report .short .intro`).First(); introShort.Length() > 0 {
		info.Intro = strings.TrimSpace(introShort.Text())
	}
	if metaCover := doc.Find(`meta[property="og:image"]`); metaCover.Length() > 0 {
		info.CoverURL, _ = metaCover.Attr("content")
	}
	if info.CoverURL == "" {
		if img := doc.Find(`#mainpic img`); img.Length() > 0 {
			info.CoverURL, _ = img.Attr("src")
		}
	}
	return info, nil
}

func (d *bookDoubanSource) setDetailHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://www.douban.com/")
	if d.cookie != "" {
		req.Header.Set("Cookie", d.cookie)
	}
}

// DownloadCover 下载封面图片
func (d *bookDoubanSource) DownloadCover(ctx context.Context, id, coverURL string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", coverURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://book.douban.com/")
	if d.cookie != "" {
		req.Header.Set("Cookie", d.cookie)
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	ext := ".jpg"
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "png") {
		ext = ".png"
	} else if strings.Contains(contentType, "webp") {
		ext = ".webp"
	}

	filename := id + ext
	out, err := CreateFile(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// CreateFile 为测试或覆盖而提取的文件创建函数（可自定义）
func CreateFile(name string) (*os.File, error) {
	return os.Create(name)
}
