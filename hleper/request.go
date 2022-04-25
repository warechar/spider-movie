package hleper

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"strings"
)

var (
	APPLICATION_JOSN = "application/json"
)

type Url struct {
	Host string
	Path string
	Query netUrl.Values
}

type Request struct {
	Url
	Header http.Header
	Client *http.Client
}

func (url *Url) check() error {
	if url.Host == "" || url.Path == "" {
		return errors.New("Host or Path is not null")
	}

	return nil
}

func (url *Url) markUrl() string {
	path := ""
	if 0 != strings.Index("/", url.Path) {
		path = "/" + url.Path
	} else {
		path = url.Path
	}

	if url.Query != nil {
		path += "?" + url.Query.Encode()
	}

	return url.Host + path
}

func NewRequest() *Request {
	return &Request{
		Url:    Url{},
		Header: make(http.Header),
		Client: nil,
	}
}

func (req *Request) Do(method string, body io.Reader) io.ReadCloser {
	if err := req.Url.check(); err != nil {
		panic(err)
	}

	newReq, err := http.NewRequest(method, req.Url.markUrl(), body)
	if err != nil {
		panic(err)
	}

	newReq.Header = req.Header

	if req.Client == nil {
		req.Client = &http.Client{}
	}

	response, err := req.Client.Do(newReq)

	if err != nil {
		panic(err)
	}

	return response.Body
}

func (req *Request) Get() string {
	body := req.Do("GET", nil)
	defer body.Close()

	dom, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return string(dom)
}

func (req *Request) Post(params map[string]interface{}) string {
	data, err := json.Marshal(params)

	if err != err{
		panic(err)
	}

	body := req.Do("POST", strings.NewReader(string(data)))
	defer body.Close()

	dom, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return string(dom)
}