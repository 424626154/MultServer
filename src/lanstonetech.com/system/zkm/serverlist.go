package zkm

import (
	"fmt"
)

const (
	SERVERLISTNODE = "serverlist"
)

type ZKServerInfo struct {
	Group     string `json:"group"`
	Type      uint8  `json:"type"`
	Index     int    `json:"index"`
	IP        string `json:"ip"`
	Port      uint16 `json:"port"`
	LocalIP   string `json:"local_ip"`
	LocalPort uint16 `json:"local_port"`
	Domain    string `json:"domain"`
}

var Server server

type server struct {
}

func (this *server) MakeRootNode() string {
	return Root()
}

func (this *server) MakeServerListNode() string {
	return fmt.Sprintf("%s/serverlist", this.MakeRootNode())
}

func (this *server) MakeServerListGroupNode(group string, class uint8) string {
	return fmt.Sprintf("%s/%s_%d", this.MakeServerListNode(), group, class)
}

func (this *server) MakeServerListItemNode(group string, class uint8, index int, ip string, port uint16, local_ip string, local_port uint16, domain string) string {
	return fmt.Sprintf(`%s/{"group":"%s","type":%d,"index":%d,"ip":"%s","port":%d,"local_ip":"%s","local_port":%d,"domain":"%s"}`, this.MakeServerListGroupNode(group, class), group, class, index, ip, port,
		local_ip, local_port, domain)
}

func (this *server) Register(group string, class uint8, index int, ip string, port uint16, local_ip string, local_port uint16, domain string) error {
	item := this.MakeServerListItemNode(group, class, index, ip, port, local_ip, local_port, domain)
	if exist, _, err := Exists(item); err == nil && exist {
		return nil
	}

	// /root
	CreateIfNotExists(this.MakeRootNode(), "", false)
	// /root/serverlist
	CreateIfNotExists(this.MakeServerListNode(), "", false)
	// /root/serverlist/group_index
	CreateIfNotExists(this.MakeServerListGroupNode(group, class), "", false)
	// /root/serverlist/group_index/{...}
	CreateIfNotExists(item, "", true)

	return nil
}
