package routers

import (
	"webmanager/controllers"
)

//==============================================================
//--------------------------路由表-------------------------------
//==============================================================
func Init() {
	Dax.AddRouter("/login", &controllers.LoginController{})
}
