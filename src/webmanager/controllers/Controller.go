package controllers

import (
	"fmt"
	"html/template"
	"lanstonetech.com/common"
	"lanstonetech.com/common/logger"
	"net/http"
	"reflect"
	"strings"
)

type Controller struct {
	Data map[interface{}]interface{}
	Tpl  string

	w http.ResponseWriter
	r *http.Request
}

func (this *Controller) Init(c interface{}, w http.ResponseWriter, r *http.Request) {
	this.Data = make(map[interface{}]interface{})

	this.w = w
	this.r = r

	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic: %v\n", err)
			http.NotFound(w, r)
		}
	}()

	urlList := common.GetUrlSplit(r.URL.Path)
	control := reflect.ValueOf(c)
	if len(urlList) == 1 {
		if r.Method == "GET" {
			control.MethodByName("Get").Call(nil)
		} else if r.Method == "POST" {
			control.MethodByName("Post").Call(nil)
		}
	} else {
		preurl := urlList[1]
		url := strings.ToUpper(preurl[:1]) + preurl[1:]

		control.MethodByName(url).Call(nil)
	}
}

func (this *Controller) ExecuteTpl() {
	t, err := template.ParseFiles(fmt.Sprintf("./views/%s", this.Tpl))
	if err != nil {
		logger.Errorf("template.ParseFiles failed! err=%v", err)
		http.NotFound(this.w, this.r)
		return
	}

	t.Execute(this.w, this.Data)
}
