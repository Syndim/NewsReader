package helpers

import (
    "io/ioutil"
    "testing"
)

func TestGetUrl(t *testing.T) {
    body, err := GetUrl("http://m.cnbeta.com/wap")
    defer body.Close()
    if err != nil {
        t.Error(err)
    } else if body == nil {
        t.Error("body == nil")
    } else {
        t.Log("Test GetUrl success")
    }
}

func TestPostUrl(t *testing.T) {
    body, err := PostUrl("http://www.cnbeta.com/cmt", "op=MSwyNzA2NTAsYTYwMWI%3DuscCHONc", &map[string]string{
        "Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
        "X-Requested-With": "XMLHttpRequest",
    })

    if err != nil {
        t.Error(err)
    } else if body == nil {
        t.Error("body == nil")
    } else {
        content, err := ioutil.ReadAll(body)
        if err != nil {
            t.Error(err)
        } else {
            t.Log(string(content))
        }
    }
}
