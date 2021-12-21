package message

// MsgID 消息 id
type MsgID uint32

// Messager 消息
type Messager interface {
	SetMsgID(MsgID)
	SetDataLength(uint32)
	SetData([]byte)

	GetDataLength() uint32
	GetMsgID() MsgID
	GetData() []byte
}

// GenMessage ...
type GenMessage func(MsgID, []byte) Messager

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
