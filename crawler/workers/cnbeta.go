package workers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"newsreader/crawler/helpers"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	MAX_PAGES          = 3
	NEWS_INSERT_SQL    = "INSERT INTO cnbeta_news(title, intro, content, origin_id, created_at, origin_site) VALUES(?, ?, ?, ?, ?, 'CnBeta')"
	NEWS_SELECT_SQL    = "SELECT origin_id FROM cnbeta_news WHERE origin_id = ?"
	COMMENT_INSERT_SQL = "INSERT INTO cnbeta_comment(origin_id, content, updated_at, origin_site) VALUES(?, ?, ?, 'CnBeta')"
	COMMENT_SELECT_SQL = "SELECT content FROM cnbeta_comment WHERE origin_id = ?"
	COMMENT_UPDATE_SQL = "UPDATE cnbeta_comment SET origin_id = ?, content = ?, updated_at = ? WHERE origin_id = ?"
)

type CommentDetails struct {
	NewsId int
	NewsSn string
	Token  string
}

type News struct {
	Title     string
	Intro     string
	Content   string
	CnbetaId  int
	CreatedAt *time.Time
}

type Comment struct {
	CnbetaId  int
	UpdatedAt *time.Time
	Content   []string
}

type Cnbeta struct {
	newsListUrlFormat  string
	newsUrlRegex       *regexp.Regexp
	tokenRegex         *regexp.Regexp
	newsUrlFormat      string
	commentUrl         string
	commentDetailRegex *regexp.Regexp
	commentPostHeaders *map[string]string
	opCodeFormat       string
}

func NewCnbeta() *Cnbeta {
	return &Cnbeta{
		newsListUrlFormat:  "http://m.cnbeta.com/wap/index.htm?page=%d",
		newsUrlRegex:       regexp.MustCompile(`/wap/view_(\d+)\.htm`),
		newsUrlFormat:      "http://www.cnbeta.com/articles/%d.htm",
		commentUrl:         "http://www.cnbeta.com/cmt",
		commentDetailRegex: regexp.MustCompile(`{SID:"(\d+?)",.*?SN:"([0-9a-fA-F]+?)"}`),
		tokenRegex:         regexp.MustCompile(`{{TOKEN:"([a-fA-F0-9]{40})"`),
		commentPostHeaders: &map[string]string{
			"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
			"X-Requested-With": "XMLHttpRequest",
		},
		opCodeFormat: "%d,%d,%s",
	}
}

func (this *Cnbeta) Start(db *sql.DB) error {
	for page := 1; page <= MAX_PAGES; page++ {
		newsIds, err := this.GetNewsList(page)
		if err != nil {
			return err
		}

		for _, newsId := range newsIds {
			fmt.Println("Start processing NewsId: " + strconv.Itoa(newsId))
			news, err := db.Query(NEWS_SELECT_SQL, newsId)
			if err != nil {
				return err
			}

			if !news.Next() {
				fmt.Println("NewsId: " + strconv.Itoa(newsId) + " not found, crawling...")
				news, err := this.GetNews(newsId)
				if err != nil {
					return err
				}

				_, err = db.Exec(
					NEWS_INSERT_SQL,
					news.Title,
					news.Intro,
					news.Content,
					news.CnbetaId,
					news.CreatedAt)

				if err != nil {
					return err
				}
			}

			comment, err := db.Query(COMMENT_SELECT_SQL, newsId)
			if err != nil {
				return err
			}

			newComment, err := this.GetAllComments(newsId)
			if err != nil {
				return err
			}

			newCommentContent := strings.Join(newComment.Content, ",")
			if comment.Next() {
				fmt.Println("Found comment for NewsId: " + strconv.Itoa(newsId))
				var commentContent string
				comment.Scan(&commentContent)
				if len(newCommentContent) > len(commentContent) {
					fmt.Println("Updating comment for NewsId: " + strconv.Itoa(newsId))
					_, err = db.Exec(
						COMMENT_UPDATE_SQL,
						newComment.CnbetaId,
						newCommentContent,
						time.Now(),
						newsId)
					if err != nil {
						return err
					}
				}
			} else {
				fmt.Println("Comment for newsId: " + strconv.Itoa(newsId) + " not found, adding...")
				_, err = db.Exec(
					COMMENT_INSERT_SQL,
					newComment.CnbetaId,
					newCommentContent,
					time.Now())
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (this *Cnbeta) GetNews(newsId int) (*News, error) {
	doc, err := goquery.NewDocument(fmt.Sprintf(this.newsUrlFormat, newsId))
	if err != nil {
		return nil, err
	}

	titleNode := doc.Find("h2#news_title")
	titleText := titleNode.Text()
	createdAtNode := doc.Find("div.title_bar span.date")
	createdAtText := createdAtNode.Text()
	createdAtTime, err := time.Parse("2006-01-02 15:04:05", createdAtText)
	introNode := doc.Find("div.introduction")
	introText := strings.TrimSpace(introNode.Text())
	contentNode := doc.Find("div.content")
	contentText, err := contentNode.Html()
	if err != nil {
		return nil, err
	}

	return &News{
		Title:     titleText,
		CreatedAt: &createdAtTime,
		Intro:     introText,
		Content:   contentText,
		CnbetaId:  newsId,
	}, nil
}

func (this *Cnbeta) GetNewsList(pageNum int) ([]int, error) {
	newsListUrl := fmt.Sprintf(this.newsListUrlFormat, pageNum)
	body, err := helpers.GetUrl(newsListUrl)
	defer body.Close()
	if err != nil {
		return nil, err
	}

	pageContent, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	newsList := this.newsUrlRegex.FindAllSubmatch(pageContent, -1)
	result := make([]int, len(newsList))

	for index, value := range newsList {
		newsId, err := strconv.ParseInt(string(value[1]), 0, 0)
		if err != nil {
			return nil, err
		}

		result[index] = int(newsId)
	}

	return result, nil
}

func (this *Cnbeta) GetCommentDetails(newsId int) (*CommentDetails, error) {
	newsUrl := fmt.Sprintf(this.newsUrlFormat, newsId)
	body, err := helpers.GetUrl(newsUrl)
	if err != nil {
		return nil, err
	}

	pageContent, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	commentDetails := this.commentDetailRegex.FindSubmatch(pageContent)
	token := this.tokenRegex.FindSubmatch(pageContent)

	return &CommentDetails{
		NewsId: newsId,
		NewsSn: string(commentDetails[2]),
		Token:  string(token[1]),
	}, nil
}

func (this *Cnbeta) GetOpCode(page int, commentDetails *CommentDetails) string {
	opCodeStr := fmt.Sprintf(
		this.opCodeFormat,
		page,
		commentDetails.NewsId,
		commentDetails.NewsSn)
	return opCodeStr
}

func (this *Cnbeta) GetComment(opCode string, token string) (io.ReadCloser, error) {
	return helpers.PostUrl(this.commentUrl, "op="+opCode+"&csrf_token="+token, this.commentPostHeaders)
}

func (this *Cnbeta) GetAllComments(newsId int) (*Comment, error) {
	commentDetails, err := this.GetCommentDetails(newsId)
	if err != nil {
		return nil, err
	}

	var result []string

	page := 1

	for {
		opCode := this.GetOpCode(page, commentDetails)
		commentBody, err := this.GetComment(opCode, commentDetails.Token)
		defer commentBody.Close()
		if err != nil {
			return nil, err
		}

		commentContent, err := ioutil.ReadAll(commentBody)
		if err != nil {
			return nil, err
		}

		var commentResult map[string]string
		err = json.Unmarshal(commentContent, &commentResult)
		if err != nil {
			return nil, err
		}

		if commentResult["status"] == "error" && commentResult["result"] == "busy" {
			time.Sleep(time.Second)
			continue
		}

		decodedCommentResult, err := base64.StdEncoding.DecodeString(commentResult["result"])
		if err != nil {
			return nil, err
		}

		var decodedCommentResultMap map[string]interface{}
		json.Unmarshal(decodedCommentResult, &decodedCommentResultMap)
		if comments, ok := decodedCommentResultMap["cmntstore"]; ok && len(comments.(map[string]interface{})) > 0 {
			result = append(result, commentResult["result"])
			page++
		} else {
			break
		}
	}

	return &Comment{
		CnbetaId: newsId,
		Content:  result,
	}, nil
}
