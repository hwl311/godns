package message

import (
	"log"
	"net"
)

const (
	CLASS_IN int = 1 //the Internet
	CLASS_CS int = 2 //the CSNET class (Obsolete - used only for examples in some obsolete RFCs)
	CLASS_CH int = 3 //the CHAOS class
	CLASS_HS int = 4 //Hesiod [Dyer 87]
)

const (
	QTYPE_A     int = 1  //a host address
	QTYPE_NS    int = 2  //an authoritative name server
	QTYPE_MD    int = 3  //a mail destination (Obsolete - use MX)
	QTYPE_MF    int = 4  //a mail forwarder (Obsolete - use MX)
	QTYPE_CNAME int = 5  //the canonical name for an alias
	QTYPE_SOA   int = 6  //marks the start of a zone of authority
	QTYPE_MB    int = 7  //a mailbox domain name (EXPERIMENTAL)
	QTYPE_MG    int = 8  //a mail group member (EXPERIMENTAL)
	QTYPE_MR    int = 9  //a mail rename domain name (EXPERIMENTAL)
	QTYPE_NULL  int = 10 //a null RR (EXPERIMENTAL)
	QTYPE_WKS   int = 11 //a well known service description
	QTYPE_PTR   int = 12 //a domain name pointer
	QTYPE_HINFO int = 13 //host information
	QTYPE_MINFO int = 14 //mailbox or mail list information
	QTYPE_MX    int = 15 //mail exchange
	QTYPE_TXT   int = 16 //text strings
	QTYPE_AAAA  int = 28 //text strings
)

const (
	CODE_OK              int = 0 //No error condition，没有错误条件；
	CODE_FORMAT_ERROR    int = 1 //Format error，请求格式有误，服务器无法解析请求；
	CODE_SERVER_FAILURE  int = 2 //Server failure，服务器出错。
	CODE_NAME_ERROR      int = 3 //Name Error，只在权威DNS服务器的响应中有意义，表示请求中的域名不存在。
	CODE_NOT_IMPLEMENTED int = 4 //Not Implemented，服务器不支持该请求类型。
	CODE_REFUSED         int = 5 //Refused，服务器拒绝执行请求操作。
)

type Message struct {
	Id     int //占16位。该值由发出DNS请求的程序生成，DNS服务器在响应时会使用该ID，这样便于请求程序区分不同的DNS响应。
	QR     int //占1位。指示该消息是请求还是响应。0表示请求；1表示响应。
	Opcode int //占4位。指示请求的类型，有请求发起者设定，响应消息中复用该值。0表示标准查询；1表示反转查询；2表示服务器状态查询。3~15目前保留，以备将来使用。
	AA     int //（Authoritative Answer，权威应答）：占1位。表示响应的服务器是否是权威DNS服务器。只在响应消息中有效。
	TC     int //（TrunCation，截断）：占1位。指示消息是否因为传输大小限制而被截断。
	RD     int //（Recursion Desired，期望递归）：占1位。该值在请求消息中被设置，响应消息复用该值。如果被设置，表示希望服务器递归查询。但服务器不一定支持递归查询。
	RA     int //（Recursion Available，递归可用性）：占1位。该值在响应消息中被设置或被清除，以表明服务器是否支持递归查询。
	Z      int //占3位。保留备用。

	RCODE int //占4位。该值在响应消息中被设置。取值及含义如下：
	//0：No error condition，没有错误条件；
	//1：Format error，请求格式有误，服务器无法解析请求；
	//2：Server failure，服务器出错。
	//3：Name Error，只在权威DNS服务器的响应中有意义，表示请求中的域名不存在。
	//4：Not Implemented，服务器不支持该请求类型。
	//5：Refused，服务器拒绝执行请求操作。
	//6~15：保留备用。

	QDCOUNT int //占16位（无符号）。指明Question部分的包含的实体数量。
	ANCOUNT int //占16位（无符号）。指明Answer部分的包含的RR（Resource Record）数量。
	NSCOUNT int //占16位（无符号）。指明Authority部分的包含的RR（Resource Record）数量。
	ARCOUNT int //占16位（无符号）。指明Additional部分的包含的RR（Resource Record）数量。

	Question   []*Question
	Answer     []*Answer
	Authority  []*Answer
	Additional []*Answer
}

func Parse(data []byte) *Message {
	m := &Message{}
	parseHeader(data, m)

	begin := 12

	qs := make([]*Question, m.QDCOUNT)
	for i := 0; i < m.QDCOUNT; i++ {
		question, size := parseQuestion(data[begin:])
		begin = begin + size
		qs[i] = question
	}
	m.Question = qs

	as := make([]*Answer, m.ANCOUNT)
	for i := 0; i < m.ANCOUNT; i++ {
		answer, size := parseAnswer(data[begin:], data)
		answer.Print()
		begin = begin + size
		as[i] = answer
	}
	m.Answer = as

	as = make([]*Answer, m.NSCOUNT)
	for i := 0; i < m.NSCOUNT; i++ {
		answer, size := parseAnswer(data[begin:], data)
		answer.Print()
		begin = begin + size
		as[i] = answer
	}
	m.Authority = as

	as = make([]*Answer, m.ARCOUNT)
	for i := 0; i < m.ARCOUNT; i++ {
		answer, size := parseAnswer(data[begin:], data)
		answer.Print()
		begin = begin + size
		as[i] = answer
	}
	m.Additional = as
	return m
}

func parseHeader(data []byte, m *Message) {
	m.Id = (int(data[0]) << 8) + int(data[1])
	m.QR = int((data[2] & 0x80) >> 7)
	m.Opcode = int((data[2] & 0x78) >> 3)
	m.AA = int((data[2] & 0x04) >> 2)
	m.TC = int((data[2] & 0x02) >> 1)
	m.RD = int((data[2] & 0x01))
	m.RA = int((data[3] & 0x80) >> 7)
	m.Z = int((data[3] & 0x70) >> 4)
	m.RCODE = int((data[3] & 0x0f))
	m.QDCOUNT = (int(data[4]) << 8) + int(data[5])
	m.ANCOUNT = (int(data[6]) << 8) + int(data[7])
	m.NSCOUNT = (int(data[8]) << 8) + int(data[9])
	m.ARCOUNT = (int(data[10]) << 8) + int(data[11])
}

type Question struct {
	QName  string
	QType  int
	QClass int
}

func (q *Question) ToBytes() []byte {
	return nil
}

func (q *Question) Name() string {
	return q.QName
}

func name(data []byte) string {
	buf := make([]byte, len(data))
	copy(buf, data)
	pos := 0
	for {
		if buf[pos] == 0 {
			break
		}
		n := int(buf[pos])
		buf[pos] = '.'
		pos = pos + n + 1
	}
	return string(buf[1:pos])
}

func parseQuestion(data []byte) (*Question, int) {
	q := &Question{}
	pos := 0
	for {
		if data[pos] == 0 {
			break
		}
		pos = pos + int(data[pos]) + 1
	}
	q.QName = name(data[:pos+1])
	q.QType = (int(data[pos+1]) << 8) + int(data[pos+2])
	q.QClass = (int(data[pos+3]) << 8) + int(data[pos+4])

	return q, pos + 5
}

type Answer struct {
	Name  string
	Data  string
	Type  int
	Class int
	Ttl   int
}

type Value interface {
	QType() int
	QName() string
	Value() string
	Data() []byte
	Format() string
}

func (a *Answer) Print() {
	log.Println("--------------------")
	log.Println("name", a.Name)
	switch a.Type {
	case QTYPE_A:
		log.Println("A", a.Data)
	case QTYPE_AAAA:
		log.Println("AAAA", a.Data)
	case QTYPE_CNAME:
		log.Println("CNAME", a.Data)
	case QTYPE_NS:
		log.Println("NS", a.Data)
	case QTYPE_TXT:
		log.Println("TXT", a.Data)
	default:
		log.Println("raw", a.Type, a.Data)

	}
}

func (a *Answer) ToBytes() []byte {
	return nil
}

func NewAnswer(Type, Class, Ttl int, value string) *Answer {
	a := &Answer{"", "", Type, Class, Ttl}
	return a
}

func parseAnswer(data []byte, raw []byte) (*Answer, int) {
	a := &Answer{}
	len := 0
	var value []byte
	if (data[0] & 0xc0) == 0xc0 {
		pos := (int(data[0]) & 0x3f) + int(data[1])
		end := 0
		for {
			if raw[pos+end] == 0 {
				break
			}
			end = end + 1
		}
		a.Name = name(raw[pos : pos+end+1])
		a.Type = (int(data[2]) << 8) + int(data[3])
		a.Class = (int(data[4]) << 8) + int(data[5])
		a.Ttl = (int(data[6]) << 24) + (int(data[7]) << 16) + (int(data[8]) << 8) + int(data[9])
		size := (int(data[10]) << 8) + int(data[11])
		value = data[12 : 12+size]
		//a.raw = data[:12+size]

		len = 12 + size

	} else {
		pos := 0
		for {
			if data[pos] == 0 {
				break
			}
			pos = pos + 1
		}
		log.Println("parse", data, data[:pos+1])
		a.Name = name(data[:pos+1])

		a.Type = (int(data[pos+1]) << 8) + int(data[pos+2])
		a.Class = (int(data[pos+3]) << 8) + int(data[pos+4])
		a.Ttl = (int(data[pos+5]) << 24) + (int(data[pos+6]) << 16) + (int(data[pos+7]) << 8) + int(data[pos+8])
		size := (int(data[pos+9]) << 8) + int(data[pos+10])
		value = data[pos+11 : pos+11+size]
		//a.raw = data[:pos+11+size]

		len = pos + size
	}

	switch a.Type {
	case QTYPE_A, QTYPE_AAAA:
		a.Data = net.IP(value).String()
	case QTYPE_CNAME, QTYPE_NS:
		a.Data = name(value)
	default:
		a.Data = string(value)
		//log.Println("raw", a.Type, a.Data)

	}

	return a, len
}

type MsgBuilder struct {
	msg   *Message
	names map[string]int
}

func NewMsgBuilder() *MsgBuilder {
	msg := &Message{}
	names := make(map[string]int)
	return &MsgBuilder{msg, names}
}

func NewResMsgBuilder(request *Message) *MsgBuilder {
	msg := &Message{}
	msg.Id = request.Id
	msg.QR = 1
	msg.Opcode = request.Opcode
	msg.AA = 0
	msg.TC = 0
	msg.RD = request.RD
	msg.RCODE = 0
	names := make(map[string]int)
	mb := &MsgBuilder{msg, names}

	//msg.Question = request.Question
	for _, q := range request.Question {
		mb.AddQuestion(q.Name(), q.QType, q.QClass)
	}
	return mb
}

func (mb *MsgBuilder) SetRecursion(ra int) *MsgBuilder {
	mb.msg.RA = ra
	return mb
}

func (mb *MsgBuilder) AddQuestion(name string, qtype int, qclass int) *MsgBuilder {
	if mb.msg.Question == nil {
		mb.msg.Question = make([]*Question, 0)
		mb.msg.QDCOUNT = 0
	}
	mb.msg.Question = append(mb.msg.Question, &Question{name, qtype, qclass})
	mb.msg.QDCOUNT += 1
	return mb
}

func (mb *MsgBuilder) AddAnswer(name string, qtype int, qclass int, ttl int, value string) *MsgBuilder {
	if mb.msg.Answer == nil {
		mb.msg.Answer = make([]*Answer, 0)
		mb.msg.ANCOUNT = 0
	}
	mb.msg.Answer = append(mb.msg.Answer, &Answer{name, value, qtype, qclass, ttl})
	mb.msg.ANCOUNT += 1
	return mb
}

func (mb *MsgBuilder) AddAuthority(name string, qtype int, qclass int, ttl int, value string) *MsgBuilder {
	if mb.msg.Authority == nil {
		mb.msg.Authority = make([]*Answer, 0)
		mb.msg.NSCOUNT = 0
	}
	mb.msg.Authority = append(mb.msg.Authority, &Answer{name, value, qtype, qclass, ttl})
	mb.msg.NSCOUNT += 1
	return mb
}

func (mb *MsgBuilder) AddAdditional(name string, qtype int, qclass int, ttl int, value string) *MsgBuilder {
	if mb.msg.Additional == nil {
		mb.msg.Additional = make([]*Answer, 0)
		mb.msg.ARCOUNT = 0
	}
	mb.msg.Additional = append(mb.msg.Additional, &Answer{name, value, qtype, qclass, ttl})
	mb.msg.ARCOUNT += 1
	return mb
}

func (mb *MsgBuilder) ToMessage() *Message {
	return mb.msg
}

func (mb *MsgBuilder) ToBytes() []byte {
	buf := make([]byte, 4096)
	msg := mb.msg
	buf[0] = byte((msg.Id >> 8) & 0xff)
	buf[1] = byte(msg.Id & 0xff)
	buf[2] = byte((msg.QR << 7) | (msg.Opcode << 3) | (msg.AA << 2) | (msg.TC << 1) | (msg.RD))
	buf[3] = byte((msg.RA << 7) | (msg.Z << 4) | (msg.RCODE))

	buf[4] = byte((msg.QDCOUNT >> 8) & 0xff)
	buf[5] = byte(msg.QDCOUNT & 0xff)

	buf[6] = byte((msg.ANCOUNT >> 8) & 0xff)
	buf[7] = byte(msg.ANCOUNT & 0xff)

	buf[8] = byte((msg.NSCOUNT >> 8) & 0xff)
	buf[9] = byte(msg.NSCOUNT & 0xff)

	buf[10] = byte((msg.ARCOUNT >> 8) & 0xff)
	buf[11] = byte(msg.ARCOUNT & 0xff)

	pos := 12
	log.Println("0ToBytes", buf[:pos])
	//question
	for _, q := range msg.Question {
		len := NameToBytes(q.QName, buf[pos:])
		pos = pos + len
		log.Println(q.QName, len)
		buf[pos] = byte((q.QType >> 8) & 0xff)
		buf[pos+1] = byte(q.QType & 0xff)
		buf[pos+2] = byte((q.QClass >> 8) & 0xff)
		buf[pos+3] = byte(q.QClass & 0xff)
		pos = pos + 4
		log.Println("1ToBytes", buf[:pos])
	}

	//answer
	for _, a := range msg.Answer {
		pos = pos + mb.answerToBytes(a, buf[pos:])
		log.Println("2ToBytes", buf[:pos])
	}

	//authority
	for _, a := range msg.Authority {
		pos = pos + mb.answerToBytes(a, buf[pos:])
		log.Println("3ToBytes", buf[:pos])
	}

	//Additional
	for _, a := range msg.Additional {
		pos = pos + mb.answerToBytes(a, buf[pos:])
		log.Println("4ToBytes", buf[:pos])
	}
	return buf[:pos]
}

func (mb *MsgBuilder) answerToBytes(a *Answer, buf []byte) int {
	pos := 0
	len := NameToBytes(a.Name, buf[pos:])
	pos = pos + len
	//type
	buf[pos] = byte((a.Type >> 8) & 0xff)
	buf[pos+1] = byte(a.Type & 0xff)
	pos = pos + 2

	//class
	buf[pos] = byte((a.Class >> 8) & 0xff)
	buf[pos+1] = byte(a.Class & 0xff)
	pos = pos + 2

	//ttl
	buf[pos] = byte((a.Ttl >> 24) & 0xff)
	buf[pos+1] = byte((a.Ttl >> 16) & 0xff)
	buf[pos+2] = byte((a.Ttl >> 8) & 0xff)
	buf[pos+3] = byte(a.Ttl & 0xff)
	pos = pos + 4

	//value
	switch a.Type {
	case QTYPE_A:
		ipaddr, _ := net.ResolveIPAddr("ip", a.Data)
		len = copy(buf[pos+2:], []byte(ipaddr.IP.To4()))
	case QTYPE_AAAA:
		ipaddr, _ := net.ResolveIPAddr("ip", a.Data)
		len = copy(buf[pos+2:], []byte(ipaddr.IP.To16()))
	case QTYPE_CNAME, QTYPE_NS:
		len = NameToBytes(a.Data, buf[pos+2:])
	default:
		len = copy(buf[pos+2:], []byte(a.Data))
	}

	//value size
	log.Println("value size", len, byte((len>>8)&0xff), byte(len&0xff))
	buf[pos] = byte((len >> 8) & 0xff)
	buf[pos+1] = byte(len & 0xff)
	return pos + len + 2
}

func NameToBytes(name string, buf []byte) int {
	pos := 0
	len := 0
	for _, c := range name {
		if c == '.' {
			buf[pos] = byte(len & 0xff)
			pos = pos + len + 1
			len = 0
		} else {
			len = len + 1
			buf[pos+len] = byte(c & 0xff)
		}
	}
	buf[pos] = byte(len & 0xff)
	return pos + len + 2
}

func (mb *MsgBuilder) ToError(code int) []byte {
	return nil
}
