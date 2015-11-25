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

func (this *LoginController) Get() {
	this.Tpl = "login.html"
	this.ExecuteTpl()
}

func (this *LoginController) Post() {
	this.r.ParseForm()
	account := this.r.FormValue("account")
	password := this.r.FormValue("password")

	fmt.Printf("account=%v password=%v\n", account, password)
	if account == "admin" && password == "admin" {
		fmt.Printf("account=%v password=%v\n", account, password)
		http.Redirect(this.w, this.r, "/", 200)
	}
}
