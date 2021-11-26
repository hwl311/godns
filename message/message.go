package message

import "log"

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
		qs[i] = parseQuestion(data[begin:])
		begin = begin + qs[i].Size()
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
	m.Id = int(data[0]) + (int(data[1]) << 8)
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
	QName  []byte
	QType  int
	QClass int
	raw    []byte
}

func (q *Question) ToBytes() []byte {
	return q.raw
}

func (q *Question) Size() int {
	return len(q.QName) + 4
}

func (q *Question) Name() string {
	return name(q.QName)
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
	return string(buf[1:])
}

func parseQuestion(data []byte) *Question {
	q := &Question{}
	pos := 0
	for {
		if data[pos] == 0 {
			break
		}
		pos = pos + int(data[pos]) + 1
	}
	q.QName = data[:pos+1]
	q.QType = (int(data[pos+1]) << 8) + int(data[pos+2])
	q.QClass = (int(data[pos+3]) << 8) + int(data[pos+4])
	q.raw = data[:pos+5]
	return q
}

type Answer struct {
	Name  []byte
	Data  []byte
	Type  int
	Class int
	Ttl   int
	raw   []byte
	AType int // 1-Answer,2-Authority or 3-Additional
}

func (a *Answer) Print() {
	log.Println("--------------------")
	log.Println("name", name(a.Name))
	switch a.Type {
	case QTYPE_A:
		log.Println("A", a.Data)
	case QTYPE_AAAA:
		log.Println("AAAA", a.Data)
	case QTYPE_CNAME:
		log.Println("CNAME", name(a.Data))
	case QTYPE_NS:
		log.Println("NS", name(a.Data))
	case QTYPE_TXT:
		log.Println("TXT", string(a.Data))
	default:
		log.Println("raw", a.Type, a.Data)

	}
}

func (a *Answer) ToBytes() []byte {
	if a.raw != nil {
		return a.raw
	}
	return nil
}

func NewAnswer(Type, Class, Ttl int, value string) *Answer {
	a := &Answer{nil, nil, Type, Class, Ttl, nil, 1}
	return a
}

func parseAnswer(data []byte, raw []byte) (*Answer, int) {
	a := &Answer{}
	if (data[0] & 0xc0) == 0xc0 {
		pos := (int(data[0]) & 0x3f) + int(data[1])
		end := 0
		for {
			if raw[pos+end] == 0 {
				break
			}
			end = end + 1
		}
		a.Name = raw[pos : pos+end+1]
		a.Type = (int(data[2]) << 8) + int(data[3])
		a.Class = (int(data[4]) << 8) + int(data[5])
		a.Ttl = (int(data[6]) << 24) + (int(data[7]) << 16) + (int(data[8]) << 8) + int(data[9])
		size := (int(data[10]) << 8) + int(data[11])
		a.Data = data[12 : 12+size]
		a.raw = data[:12+size]

		return a, 12 + size

	} else {
		pos := 0
		for {
			if data[pos] == 0 {
				break
			}
			pos = pos + 1
		}
		a.Name = data[:pos+1]

		a.Type = (int(data[pos+1]) << 8) + int(data[pos+2])
		a.Class = (int(data[pos+3]) << 8) + int(data[pos+4])
		a.Ttl = (int(data[pos+5]) << 24) + (int(data[pos+6]) << 16) + (int(data[pos+7]) << 8) + int(data[pos+8])
		size := (int(data[pos+9]) << 8) + int(data[pos+10])
		a.Data = data[pos+11 : pos+11+size]
		a.raw = data[:pos+11+size]

		return a, pos + size
	}
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

func (mb *MsgBuilder) SetRecursion(ra bool) *MsgBuilder {
	if ra {
		mb.msg.RA = 1
	} else {
		mb.msg.RA = 0
	}
	return mb
}

func (mb *MsgBuilder) AddQuestion(name string, qtype int, qclass int) *MsgBuilder {
	if mb.msg.Question == nil {
		mb.msg.Question = make([]*Question, 0)
	}
	mb.msg.Question = append(mb.msg.Question, &Question{})
	return mb
}

func (mb *MsgBuilder) AddAnswer(AnswerType int, name string, qtype int, qclass int, ttl int, value []byte) *MsgBuilder {
	if mb.msg.Answer == nil {
		mb.msg.Answer = make([]*Answer, 0)
	}
	mb.msg.Answer = append(mb.msg.Answer, &Answer{[]byte(name), value, qtype, qclass, ttl, nil, AnswerType})
	return mb
}

func (mb *MsgBuilder) ToMessage() *Message {
	return mb.msg
}

func (mb *MsgBuilder) ToBytes() []byte {
	return nil
}

func (mb *MsgBuilder) ToError(code int) []byte {
	return nil
}
