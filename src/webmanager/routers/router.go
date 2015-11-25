package routers

import (
	"webmanager/controllers"
)

func Init() {
	Dax.AddRouter("/", &controllers.MainController{})
	Dax.AddRouter("/login", &controllers.LoginController{})
}
