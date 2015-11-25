package controllers

import (
	"fmt"
	"net/http"
)

type MainController struct {
	Controller
}

func (this *MainController) Handler(w http.ResponseWriter, r *http.Request) {
	this.Controller.Init(this, w, r)
}

func (this *MainController) Login() {
	if this.r.Method == "GET" {
		fmt.Fprintf(this.w, "Login--------/login--------------Get")
		// this.Tpl = "login.html"
		// this.ExecuteTpl()
	} else if this.r.Method == "POST" {
		fmt.Fprintf(this.w, "Login--------/login--------------Post")
	}
}

func (this *MainController) Get() {
	this.Tpl = "database.html"
	this.ExecuteTpl()
}

func (this *MainController) Post() {
	fmt.Fprintf(this.w, "----------------/------Post")
}
