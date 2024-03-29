package main

import (
	"common/assert"
	"common/config"
	"io/ioutil"
	"net/http"

	"common/initializer"

	"github.com/Sirupsen/logrus"
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	incData, err := ioutil.ReadAll(r.Body)
	assert.Nil(err)
	defer func() {
		err := r.Body.Close()
		assert.Nil(err)
	}()

	req, err := http.NewRequest(r.Method, config.Config.Proxy.URL+r.RequestURI, r.Body)
	assert.Nil(err)

	var resp *http.Response
	resp, err = sendReq(req)
	assert.Nil(err)

	var respData []byte
	respData, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	logrus.WithFields(logrus.Fields{
		"request":  string(incData),
		"response": string(respData),
	}).Info(r.URL)

}

func sendReq(req *http.Request) (*http.Response, error) {
	Client := &http.Client{}
	return Client.Do(req)

}

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()

	assert.Nil(http.ListenAndServe(config.Config.Proxy.Port, &handler{}))
}
