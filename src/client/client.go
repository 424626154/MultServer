package main

import (
	"fmt"
	"lanstonetech.com/common"
	"lanstonetech.com/network"
	"net"
)

func main() {
	common_package := new(network.CommonPackage)
	common_package.Account = "dcea85c06887c52f20994c393c5bbcdc"
	common_package.Signature = "dcea85c06887c52f20994c393c5bbcdc"
	common_package.Token = "dcea85c06887c52f20994c393c5bbcdc"
	greeting := "Hello LoginServer"

	size := uint32(4 + 4 + 2 + 48 + 2 + 48 + 2 + 48 + 2 + len(greeting))
	header := new(network.Header)
	header.PacketID = 100
	header.PackageLen = size - 8

	send_data := make([]byte, size)

	pos := 0
	common.WriteUint32(send_data[pos:pos+4], header.PacketID)
	pos += 4
	common.WriteUint32(send_data[pos:pos+4], header.PackageLen)
	pos += 4

	leng := len(common_package.Account)
	common.WriteUint16(send_data[pos:pos+2], uint16(leng))
	pos += 2
	common.WriteString(send_data[pos:pos+leng], common_package.Account)
	pos += 48

	leng = len(common_package.Signature)
	common.WriteUint16(send_data[pos:pos+2], uint16(leng))
	pos += 2
	common.WriteString(send_data[pos:pos+leng], common_package.Signature)
	pos += 48

	leng = len(common_package.Token)
	common.WriteUint16(send_data[pos:pos+2], uint16(leng))
	pos += 2
	common.WriteString(send_data[pos:pos+leng], common_package.Token)
	pos += 48

	leng = len(greeting)
	fmt.Printf("GreetingLen=%v\n", leng)
	common.WriteUint16(send_data[pos:pos+2], uint16(leng))
	pos += 2
	common.WriteString(send_data[pos:pos+leng], greeting)
	pos += 48

	tcpaddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Printf("[client] net.ResolveTCPAddr failed! err=%v\n", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		fmt.Printf("[client] net.DiaTCP failed! err=%v\n", err)
	}
	defer func() {
		conn.Close()
		fmt.Printf("[client] ok...!\n")
	}()

	leng, err = conn.Write(send_data)
	if err != nil {
		fmt.Printf("[client] conn.Write failed! err=%v\n", err)
	}

	fmt.Printf("[client] conn.Write len=%v\n", leng)

	head := make([]byte, 8)
	leng, err = conn.Read(head)
	if err != nil {
		fmt.Printf("conn.Read failed! err=%v\n", err)
		return
	}
	fmt.Printf("read leng=%v\n", leng)

	package_id := common.ReadUint32(head[0:4])
	package_len := common.ReadUint32(head[4:])
	fmt.Printf("package_id=%v package_len=%v\n", package_id, package_len)

	body := make([]byte, package_len)
	leng, err = conn.Read(body)
	if err != nil {
		fmt.Printf("conn.Read failed! err=%v\n", err)
		return
	}
	fmt.Printf("read leng=%v\n", leng)
	result := common.ReadUint16(head[0:2])
	fmt.Printf("result=%v\n", result)

	return
}
