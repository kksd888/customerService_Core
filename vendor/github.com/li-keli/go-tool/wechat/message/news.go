package message

//News 图文消息
type News struct {
	CommonToken

	ArticleCount int        `json:"article_count" xml:"ArticleCount"`
	Articles     []*Article `json:"articles" xml:"Articles>item,omitempty"`
}

//NewNews 初始化图文消息
func NewNews(articles []*Article) *News {
	news := new(News)
	news.ArticleCount = len(articles)
	news.Articles = articles
	return news
}

//Article 单篇文章
type Article struct {
	Title       string `json:"title" xml:"Title,omitempty"`
	Description string `json:"description" xml:"Description,omitempty"`
	PicURL      string `json:"pic_url" xml:"PicUrl,omitempty"`
	URL         string `json:"url" xml:"Url,omitempty"`
}

//NewArticle 初始化文章
func NewArticle(title, description, picURL, url string) *Article {
	article := new(Article)
	article.Title = title
	article.Description = description
	article.PicURL = picURL
	article.URL = url
	return article
}
