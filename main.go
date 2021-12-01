package main

import (
	"encoding/json"
	"log"

	"ruotian.vip/godns/message"
	"ruotian.vip/godns/server"
)

type Backend struct{}

func (be *Backend) Query(msg *message.Message, q *message.Question, mb *message.MsgBuilder) error {
	log.Println("QUESTION", q.Name(), q.QType)
	mb.AddAnswer(string(q.QName), q.QType, q.QClass, 600, "8.8.8.8")
	return nil
}
func (be *Backend) STATUS() int {
	return 0
} //服务器状态
func (be *Backend) RecursionAvailable() int {
	return 1
} //服务器是否支持递归查询

func main() {
	buf := make([]byte, 40)
	len := message.NameToBytes("www.baidu.com", buf)
	log.Println(buf[:len])
	//Test()
	server.New(&Backend{}).ListenAndServe("127.0.0.1", 54)

}
func Test() {
	log.Println("dns started!")
	parse([]byte("\xbbh\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\x03www\x07ruotian\x03vip\x00\x00\x01\x00\x01"))
	log.Println("dns started!")
	parse([]byte("\xbbh\x81\x80\x00\x01\x00\x01\x00\x00\x00\x00\x03www\x07ruotian\x03vip\x00\x00\x01\x00\x01\xc0\x0c\x00\x01\x00\x01\x00\x00\x02X\x00\x04j47\x90"))
	log.Println("dns started!")
	parse([]byte("\xbbh\x81\x80\x00\x01\x00\x02\x00\x00\x00\x00\x03www\x05baidu\x03vip\x00\x00\x01\x00\x01\xc0\x0c\x00\x01\x00\x01\x00\x00\x02X\x00\x04\x05\x16\x91\x10\xc0\x0c\x00\x01\x00\x01\x00\x00\x02X\x00\x04\x05\x16\x91y"))
	parse([]byte{187, 104, 133, 0, 0, 1, 0, 2, 0, 5, 0, 10, 3, 119, 119, 119, 5, 98, 97, 105, 100, 117, 3, 118, 105, 112, 0, 0, 1, 0, 1, 192, 12, 0, 1, 0, 1, 0, 0, 2, 88, 0, 4, 5, 22, 145, 16, 192, 12, 0, 1, 0, 1, 0, 0, 2, 88, 0, 4, 5, 22, 145, 121, 192, 16, 0, 2, 0, 1, 0, 0, 1, 44, 0, 21, 3, 110, 115, 50, 12, 98, 114, 97, 110, 100, 115, 104, 101, 108, 116, 101, 114, 2, 100, 101, 0, 192, 16, 0, 2, 0, 1, 0, 0, 1, 44, 0, 22, 3, 110, 115, 49, 12, 98, 114, 97, 110, 100, 115, 104, 101, 108, 116, 101, 114, 3, 99, 111, 109, 0, 192, 16, 0, 2, 0, 1, 0, 0, 1, 44, 0, 21, 3, 110, 115, 53, 12, 98, 114, 97, 110, 100, 115, 104, 101, 108, 116, 101, 114, 2, 117, 115, 0, 192, 16, 0, 2, 0, 1, 0, 0, 1, 44, 0, 22, 3, 110, 115, 52, 12, 98, 114, 97, 110, 100, 115, 104, 101, 108, 116, 101, 114, 3, 110, 101, 116, 0, 192, 16, 0, 2, 0, 1, 0, 0, 1, 44, 0, 23, 3, 110, 115, 51, 12, 98, 114, 97, 110, 100, 115, 104, 101, 108, 116, 101, 114, 4, 105, 110, 102, 111, 0, 192, 108, 0, 1, 0, 1, 0, 0, 1, 44, 0, 4, 178, 33, 33, 199, 192, 108, 0, 28, 0, 1, 0, 0, 14, 16, 0, 16, 32, 1, 65, 208, 0, 12, 3, 136, 1, 120, 0, 51, 0, 51, 1, 153, 192, 75, 0, 1, 0, 1, 0, 0, 168, 192, 0, 4, 66, 206, 2, 254, 192, 75, 0, 28, 0, 1, 0, 0, 168, 192, 0, 16, 38, 4, 69, 0, 0, 12, 0, 3, 0, 102, 2, 6, 0, 2, 2, 84, 192, 209, 0, 1, 0, 1, 0, 0, 84, 96, 0, 4, 94, 23, 156, 156, 192, 209, 0, 28, 0, 1, 0, 0, 168, 192, 0, 16, 32, 1, 65, 208, 0, 12, 3, 136, 0, 148, 0, 35, 1, 86, 1, 86, 192, 175, 0, 1, 0, 1, 0, 0, 14, 16, 0, 4, 192, 95, 19, 131, 192, 175, 0, 28, 0, 1, 0, 0, 14, 16, 0, 16, 38, 7, 83, 0, 0, 96, 94, 28, 1, 146, 0, 149, 0, 25, 1, 49, 192, 142, 0, 1, 0, 1, 0, 0, 14, 16, 0, 4, 192, 95, 17, 188, 192, 142, 0, 28, 0, 1, 0, 0, 14, 16, 0, 16, 38, 7, 83, 0, 0, 96, 107, 190, 1, 146, 0, 149, 0, 23, 1, 136})
	parse([]byte("\xbbh\x85\x00\x00\x01\x00\x02\x00\x05\x00\n\x03www\x05baidu\x03vip\x00\x00\x01\x00\x01\xc0\x0c\x00\x01\x00\x01\x00\x00\x02X\x00\x04\x05\x16\x91\x10\xc0\x0c\x00\x01\x00\x01\x00\x00\x02X\x00\x04\x05\x16\x91y\xc0\x10\x00\x02\x00\x01\x00\x00\x01,\x00\x15\x03ns2\x0cbrandshelter\x02de\x00\xc0\x10\x00\x02\x00\x01\x00\x00\x01,\x00\x16\x03ns4\x0cbrandshelter\x03net\x00\xc0\x10\x00\x02\x00\x01\x00\x00\x01,\x00\x17\x03ns3\x0cbrandshelter\x04info\x00\xc0\x10\x00\x02\x00\x01\x00\x00\x01,\x00\x15\x03ns5\x0cbrandshelter\x02us\x00\xc0\x10\x00\x02\x00\x01\x00\x00\x01,\x00\x16\x03ns1\x0cbrandshelter\x03com\x00\xc0\xd2\x00\x01\x00\x01\x00\x00\x01,\x00\x04\xb2!!\xc7\xc0\xd2\x00\x1c\x00\x01\x00\x00\x0e\x10\x00\x10 \x01A\xd0\x00\x0c\x03\x88\x01x\x003\x003\x01\x99\xc0K\x00\x01\x00\x01\x00\x00\xa8\xc0\x00\x04B\xce\x02\xfe\xc0K\x00\x1c\x00\x01\x00\x00\xa8\xc0\x00\x10&\x04E\x00\x00\x0c\x00\x03\x00f\x02\x06\x00\x02\x02T\xc0\x8e\x00\x01\x00\x01\x00\x00T`\x00\x04^\x17\x9c\x9c\xc0\x8e\x00\x1c\x00\x01\x00\x00\xa8\xc0\x00\x10 \x01A\xd0\x00\x0c\x03\x88\x00\x94\x00#\x01V\x01V\xc0l\x00\x01\x00\x01\x00\x00\x0e\x10\x00\x04\xc0_\x13\x83\xc0l\x00\x1c\x00\x01\x00\x00\x0e\x10\x00\x10&\x07S\x00\x00`^\x1c\x01\x92\x00\x95\x00\x19\x011\xc0\xb1\x00\x01\x00\x01\x00\x00\x0e\x10\x00\x04\xc0_\x11\xbc\xc0\xb1\x00\x1c\x00\x01\x00\x00\x0e\x10\x00\x10&\x07S\x00\x00`k\xbe\x01\x92\x00\x95\x00\x17\x01\x88"))
}

func parse(d []byte) {
	m := message.Parse(d)
	j, _ := json.MarshalIndent(m, "", "    ")
	log.Println(string(j))
}
