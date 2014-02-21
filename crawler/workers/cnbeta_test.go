package workers

import (
    "io/ioutil"
    "strconv"
    "strings"
    "testing"
)

func TestGetNews(t *testing.T) {
    cnbeta := NewCnbeta()
    news, err := cnbeta.GetNews(270917)
    if err != nil {
        t.Error(err)
    } else if news == nil {
        t.Error("news is nil")
    } else {
        if news.Title != "[图]Windows 8和8.1的全球系统市场占比达到了10.58%" {
            t.Error("news load failed, expected: [图]Windows 8和8.1的全球系统市场占比达到了10.58%, get: " + news.Title)
        }
        if news.Content == "" {
            t.Error("news content is empty")
        }
    }
}

func TestGetNewsList(t *testing.T) {
    cnbeta := NewCnbeta()
    newsList, err := cnbeta.GetNewsList(1)
    if err != nil {
        t.Error(err)
    } else if newsList == nil {
        t.Error("news list is nil")
    } else if len(newsList) != 30 {
        t.Error("news list length mismatch, expected: 30, get: " + strconv.Itoa(len(newsList)))
    }
}

func TestGetCommentDetails(t *testing.T) {
    cnbeta := NewCnbeta()
    commentDetails, err := cnbeta.GetCommentDetails(270917)
    if err != nil {
        t.Error(err)
    } else if commentDetails == nil {
        t.Error("comment details is nil")
    } else {
        if commentDetails.NewsId != 270917 {
            t.Error("comment details id mismatch, expect: 270917, get: " + strconv.Itoa(commentDetails.NewsId))
        }
        if commentDetails.NewsSn != "3b2c5" {
            t.Error("comment details sn mismatch, expect: 3b2c5, get: " + commentDetails.NewsSn)
        }
    }
}

func TestGetOpCode(t *testing.T) {
    cnbeta := NewCnbeta()
    opcode := cnbeta.GetOpCode(1, &CommentDetails{
        NewsId: 270917,
        NewsSn: "3b2c5",
    })
    if opcode != "MSwyNzA5MTcsM2IyYzU=" {
        t.Error("opcode mismatch, expected: MSwyNzA5MTcsM2IyYzU= get: " + opcode)
    }
}

func TestGetComment(t *testing.T) {
    cnbeta := NewCnbeta()
    body, err := cnbeta.GetComment("MSwyNzA5MTcsM2IyYzU=")
    defer body.Close()
    if err != nil {
        t.Error(err)
    } else if body == nil {
        t.Error("body is nil")
    } else {
        content, err := ioutil.ReadAll(body)
        if err != nil {
            t.Error(err)
        } else {
            contentStr := string(content)
            if !strings.Contains(contentStr, "success") {
                t.Error("get comment failed, content is:" + contentStr)
            }
        }
    }
}

func TestGetAllComments(t *testing.T) {
    cnbeta := NewCnbeta()
    result, err := cnbeta.GetAllComments(270954)
    if err != nil {
        t.Error(err)
    } else if result == nil {
        t.Log("result is empty, please use another news id")
    }
}
