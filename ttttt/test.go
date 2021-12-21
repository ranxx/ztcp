package ttttt

type Message struct {
	ID int64
}

func (m Message) GetID() int64 {
	return m.ID
}
