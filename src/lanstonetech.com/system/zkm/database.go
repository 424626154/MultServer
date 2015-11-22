package zkm

import (
	"fmt"
)

type DBServerInfo struct {
	Group    string `json:"group"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	DBName   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
	Slaveof  string `json:"slaveof"`
	Offline  int    `json:"offline"`
}

var Database database

type database struct {
}

func (this *database) MakeRoot() string {
	return Root()
}

func (this *database) MakeDatabaseRoot() string {
	return fmt.Sprintf("%s/database", this.MakeRoot())
}

func (this *database) MakeDatabase(driver string) string {
	return fmt.Sprintf("%s/%s", this.MakeDatabaseRoot(), driver)
}

func (this *database) MakeDatabaseGroup(driver, group string) string {
	return fmt.Sprintf("%s/%s", this.MakeDatabase(driver), group)
}

func (this *database) MakeSQLNode(driver, group, ip string, port int, dbname, user, password string, slaveof string, offline int) string {
	return fmt.Sprintf(`%s/{"group":"%s","ip":"%s","port":%d,"dbname":"%s","user":"%s","password":"%s",
		"slaveof":"%s","offline"%d}`, this.MakeDatabaseGroup(driver, group), group, ip, port, dbname,
		user, password, slaveof, offline)
}
