package controllers

import (
	"fmt"
	"net/http"
)

type LoginController struct {
	Controller
}

func (this *LoginController) Handler(w http.ResponseWriter, r *http.Request) {
	this.Controller.Init(this, w, r)
}

func (this *LoginController) Login() {
	this.Tpl = "login.html"
	this.ExecuteTpl()
}

func (this *LoginController) Get() {
	fmt.Fprintf(this.w, "----------------------Get")
}

func (this *LoginController) Post() {
	fmt.Fprintf(this.w, "----------------------Post")
}
