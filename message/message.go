package message

// MsgID ...
type MsgID uint32

// Messager ...
type Messager interface {
	SetMsgID(MsgID)
	SetDataLength(uint32)
	SetData([]byte)

	GetDataLength() uint32
	GetMsgID() MsgID
	GetData() []byte
}

type messager struct {
	msgid  MsgID  // 消息号
	length uint32 // 数据部分的 长度
	msg    []byte // 数据部分
}

// Empty empty
func Empty() Messager {
	return &messager{}
}

// DefaultMessager message
func DefaultMessager(msgid MsgID, msg []byte) Messager {
	return &messager{
		msgid:  msgid,
		length: uint32(len(msg)),
		msg:    msg,
	}
}

func (m *messager) SetMsgID(id MsgID) {
	m.msgid = id
}

func (m *messager) SetDataLength(length uint32) {
	m.length = length
}

func (m *messager) SetData(b []byte) {
	m.msg = b
}

func (m *messager) GetMsgID() MsgID {
	return m.msgid
}

func (m *messager) GetDataLength() uint32 {
	return m.length
}

func (m *messager) GetData() []byte {
	return m.msg
}

// // SplitFunc split
// func SplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
// 	// log.Println(atEOF, len(data))
// 	pkg := Package{}
// 	if atEOF || len(data) < int(HeadLength()) {
// 		return
// 	}

// 	if err = pkg.UnpackHeadBytes(data); err != nil {
// 		return
// 	}

// 	totalLen := int64(pkg.Length) + HeadLength()
// 	if totalLen > int64(len(data)) {
// 		return
// 	}
// 	return int(totalLen), data[:totalLen], nil
// }

// // SplitFuncEdgeTriggered ...
// func SplitFuncEdgeTriggered(data []byte, atEOF bool) (advance int, token []byte, err error) {
// 	// log.Println("client.tcp", atEOF, len(data))
// 	if atEOF && len(data) == 0 {
// 		return 0, nil, nil
// 	}
// 	if !atEOF {
// 		return len(data), data[:], nil
// 	}
// 	return 0, nil, nil
// }

// // NewScanner ...
// func NewScanner(reader io.Reader, split bufio.SplitFunc) *bufio.Scanner {
// 	scanner := bufio.NewScanner(reader)
// 	scanner.Split(split)
// 	scanner.Buffer(nil, 1024*1024*1024)
// 	return scanner
// }
