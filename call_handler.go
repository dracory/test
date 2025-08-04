package test

import (
	"net/http"
	"net/http/httptest"
)

func CallEndpoint(method string, f func(w http.ResponseWriter, r *http.Request), options NewRequestOptions) (body string, response *http.Response, err error) {
	req, err := NewRequest(method, "/", options)

	if err != nil {
		return "", nil, err
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(f)
	handler.ServeHTTP(recorder, req)
	body = recorder.Body.String()

	return body, recorder.Result(), nil
}

func CallStringEndpoint(method string, f func(w http.ResponseWriter, r *http.Request) string, options NewRequestOptions) (body string, response *http.Response, err error) {
	req, err := NewRequest(method, "/", options)

	if err != nil {
		return "", nil, err
	}

	recorder := httptest.NewRecorder()
	handler := StringHandler(f)
	handler.ServeHTTP(recorder, req)
	body = recorder.Body.String()

	return body, recorder.Result(), nil
}

type StringHandler func(w http.ResponseWriter, r *http.Request) string

func (h StringHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	str := h(w, r)
	StringResponse(w, r, str)
}

// StringResponse - responds with the string body
func StringResponse(w http.ResponseWriter, r *http.Request, body string) {
	contentType := w.Header().Get("Content-Type")

	if contentType == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}

	w.Write([]byte(body))
}
