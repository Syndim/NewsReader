package helpers

import (
    "errors"
    "io"
    "net/http"
    "strings"
)

func GetUrl(url string) (io.ReadCloser, error) {
    if url == "" {
        return nil, errors.New("Url is empty")
    }

    resp, err := http.Get(url)
    return getBody(resp, err)
}

func PostUrl(url string, data string, headers *map[string]string) (io.ReadCloser, error) {
    if url == "" {
        return nil, errors.New("Url is empty")
    }

    dataReader := strings.NewReader(data)
    req, err := http.NewRequest("POST", url, dataReader)
    if err != nil {
        return nil, err
    }

    for k, v := range *headers {
        req.Header.Add(k, v)
    }

    var client http.Client
    resp, err := client.Do(req)
    return getBody(resp, err)
}

func getBody(resp *http.Response, err error) (io.ReadCloser, error) {
    if err != nil {
        return nil, err
    } else if resp.StatusCode != 200 {
        return nil, errors.New("Status:" + resp.Status)
    } else {
        return resp.Body, nil
    }
}
