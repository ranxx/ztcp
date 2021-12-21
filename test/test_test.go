package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/ranxx/ztcp/conn"
	"github.com/ranxx/ztcp/pkg/dispatch"
	"github.com/ranxx/ztcp/pkg/encoding"
	"github.com/ranxx/ztcp/pkg/io/write"
	messagerr "github.com/ranxx/ztcp/pkg/message"
	"github.com/ranxx/ztcp/pkg/pack"
	"github.com/ranxx/ztcp/router"
	ttt "github.com/ranxx/ztcp/ttttt"
	"github.com/ranxx/ztcp/ttttt/ttttt"
)

func Test01(t *testing.T) {
	nnew := func(v interface{}) (interface{}, error) {
		t := reflect.TypeOf(v)
		return reflect.New(t).Interface(), nil
	}

	type ttt struct {
		Name string `json:"name"`
	}

	tt := ttt{Name: "axing"}

	fmt.Println(tt)
	fmt.Println(nnew(tt))
	fmt.Println(nnew(&tt))
}

func Test02(t *testing.T) {
	mmm := map[interface{}]string{}
	mmm["sss"] = "sss"
	mmm[1] = "1"
	fmt.Println(mmm)
}

type messager interface {
	GetID() int64
}
type Message struct {
	ID int64
}

func (m Message) GetID() int64 {
	return m.ID
}

func test03(msg messager) {
	// {
	// 	if
	// 	fmt.Println(reflect.TypeOf(msg).Kind())
	// 	return
	// }
	{
		fmt.Println(reflect.ValueOf(msg).MethodByName("GetID").Call(nil)[0].Int())
	}
	t := reflect.TypeOf(msg)

	fmt.Println("name", t.String(), t.PkgPath(), t.Name())

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		fmt.Println("name", t.String(), t.PkgPath(), t.Name())
	}
	v := reflect.New(t)
	vv, ok := v.Interface().(messager)
	fmt.Println(ok)
	vv.GetID()

	fmt.Println(v.MethodByName("GetID").Call(nil)[0].Int())
}

func Test03(t *testing.T) {
	msg := Message{ID: 4}
	test03(&msg)
	test03(msg)
	fmt.Println("-------------------------")
	msg2 := ttttt.Message{ID: 5}
	test03(&msg2)
	test03(msg2)
	fmt.Println("-------------------------")
	msg3 := ttt.Message{ID: 5}
	test03(&msg3)
	test03(msg3)
}

func Test04(t *testing.T) {
	mmm := conn.NewManager()
	route1 := router.NewRouter(messagerr.MsgID(1), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		s, ok := i.([]byte)
		log.Println("Hhhhhhh", string(s), ok)
	}))
	route2 := router.NewRouter(messagerr.MsgID(2), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		s, ok := i.([]byte)
		log.Println("msg-2", string(s), ok)
	}))
	route3 := router.NewRouter(messagerr.MsgID(3), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		s, ok := i.([]byte)
		log.Println("msg-3", mmm.Get(1).RemoteAddr(), string(s), ok)
	}))
	root := router.NewRoot().AddRouter(route1, route2, route3)

	list, err := net.Listen("tcp", ":12351")
	if err != nil {
		panic(err)
	}
	idex := int64(0)

	go func() {
		c, err := net.Dial("tcp", ":12351")
		if err != nil {
			panic(err)
		}
		i := 100
		for {
			time.Sleep(time.Second)
			i++
			data, err := pack.DefaultPack(messagerr.DefaultMessager(messagerr.MsgID(1), []byte(fmt.Sprintf("我再说 %d", i))))
			if err != nil {
				panic(err)
			}
			fmt.Println(c.Write(data))
			data2, err := pack.DefaultPack(messagerr.DefaultMessager(messagerr.MsgID(2), []byte(fmt.Sprintf("我说你 %d", i))))
			if err != nil {
				panic(err)
			}
			fmt.Println(c.Write(data2))
			data2, err = pack.DefaultPack(messagerr.DefaultMessager(messagerr.MsgID(3), []byte(fmt.Sprintf("我说你是个大笨蛋 %d", i))))
			if err != nil {
				panic(err)
			}
			fmt.Println(c.Write(data2))
		}
	}()
	for {
		idex++
		c, err := list.Accept()
		if err != nil {
			panic(err)
		}
		cc := conn.NewConn(idex, c, conn.WithDispatcher(dispatch.DefaultDispatcher(root, nil)))
		mmm.AddConn(cc)
		cc.Start()
	}
}

func Test05(t *testing.T) {
	mmm := conn.NewManager()
	route1 := router.NewRouter(messagerr.MsgID(1), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		s, ok := i.([]byte)
		log.Println("Hhhhhhh", string(s), ok)
	}))
	route2 := router.NewRouter(messagerr.MsgID(2), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		s, ok := i.([]byte)
		log.Println("msg-2", string(s), ok)
	}))
	route3 := router.NewRouter(messagerr.MsgID(3), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		s, ok := i.([]byte)
		log.Println("msg-3", mmm.Get(1).RemoteAddr(), string(s), ok)
	}))
	route4 := router.NewRouter(messagerr.MsgID(4), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		s, ok := i.([]byte)
		log.Println("msg-4", mmm.Get(1).RemoteAddr(), string(s), ok)
	}))
	root := router.NewRoot().AddRouter(route1, route2, route3, route4)

	list, err := net.Listen("tcp", ":12351")
	if err != nil {
		panic(err)
	}
	idex := int64(0)

	go func() {
		c, err := net.Dial("tcp", ":12351")
		if err != nil {
			panic(err)
		}
		cc := conn.NewConn(0, c, conn.WithWriter(write.DefaultWriter(nil, write.WithTypeMsgID(func() (interface{}, messagerr.MsgID) {
			return []byte{}, 4
		}))))
		writer := cc.Writer()
		i := 100
		for {
			time.Sleep(time.Second)
			i++

			data, err := pack.DefaultPack(messagerr.DefaultMessager(messagerr.MsgID(1), []byte(fmt.Sprintf("我再说 %d", i))))
			if err != nil {
				panic(err)
			}
			fmt.Println(c.Write(data))
			data2, err := pack.DefaultPack(messagerr.DefaultMessager(messagerr.MsgID(2), []byte(fmt.Sprintf("我说你 %d", i))))
			if err != nil {
				panic(err)
			}
			fmt.Println(c.Write(data2))
			data2, err = pack.DefaultPack(messagerr.DefaultMessager(messagerr.MsgID(3), []byte(fmt.Sprintf("我说你是个大笨蛋 %d", i))))
			if err != nil {
				panic(err)
			}
			fmt.Println(c.Write(data2))
			fmt.Println(writer.WriteValue([]byte(fmt.Sprintf("我在测第4个 %d  WriteValue", i))))
			fmt.Println(writer.Write(messagerr.MsgID(4), []byte(fmt.Sprintf("我在测第4个 %d  Write", i))))
			fmt.Println(writer.WriteMessager(messagerr.DefaultMessager(messagerr.MsgID(4), []byte(fmt.Sprintf("我在测第4个 %d  WriteMessager", i)))))
		}
	}()
	for {
		idex++
		c, err := list.Accept()
		if err != nil {
			panic(err)
		}
		cc := conn.NewConn(idex, c, conn.WithDispatcher(dispatch.DefaultDispatcher(root, nil)))
		mmm.AddConn(cc)
		cc.Start()
	}
}

func TestStruct(t *testing.T) {
	mmm := conn.NewManager()
	type msgType struct {
		UserName string `json:"user_name"`
	}
	route1 := router.NewRouter(messagerr.MsgID(1), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		s, ok := i.([]byte)
		log.Println("Hhhhhhh", string(s), ok)
	}))
	route2 := router.NewRouter(messagerr.MsgID(2), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		s, ok := i.([]byte)
		log.Println("msg-2", string(s), ok)
	}))
	route3 := router.NewRouter(messagerr.MsgID(3), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		s, ok := i.([]byte)
		log.Println("msg-3", mmm.Get(1).RemoteAddr(), string(s), ok)
	}))
	route4 := router.NewRouter(messagerr.MsgID(4), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		s, ok := i.([]byte)
		log.Println("msg-4", mmm.Get(1).RemoteAddr(), string(s), ok)
	}))
	route5 := router.NewRouter(messagerr.MsgID(5), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		log.Println("msg-5", mmm.Get(1).RemoteAddr(), reflect.TypeOf(i).Kind(), fmt.Sprintf("%#v", i))
	}))
	route6 := router.NewRouter(messagerr.MsgID(6), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		log.Println("msg-6", mmm.Get(1).RemoteAddr(), reflect.TypeOf(i).Kind(), fmt.Sprintf("%#v", i))
	}))
	route7 := router.NewRouter(messagerr.MsgID(7), router.WrapHandler(func(c1 context.Context, c2 net.Conn, i interface{}) {
		// TODO: 是否在这里进行反序列化比较好
		log.Println("msg-7", mmm.Get(1).RemoteAddr(), reflect.TypeOf(i).Kind(), fmt.Sprintf("%#v", i))
	}))
	root := router.NewRoot().AddRouter(route1, route2, route3, route4, route5, route6, route7)

	list, err := net.Listen("tcp", ":12351")
	if err != nil {
		panic(err)
	}
	idex := int64(0)

	go func() {
		c, err := net.Dial("tcp", ":12351")
		if err != nil {
			panic(err)
		}
		cc := conn.NewConn(0, c, conn.WithWriter(write.DefaultWriter(nil, write.WithTypeMsgID(func() (interface{}, messagerr.MsgID) {
			return []byte{}, 4
		}), write.WithTypeMsgID(func() (interface{}, messagerr.MsgID) {
			return &msgType{}, 5
		}, func() (interface{}, messagerr.MsgID) {
			return &ttt.Message{}, 6
		}, func() (interface{}, messagerr.MsgID) {
			return &ttttt.Message{}, 7
		}))))
		writer := cc.Writer()
		i := 100
		for {
			time.Sleep(time.Second)
			i++
			data, err := pack.DefaultPack(messagerr.DefaultMessager(messagerr.MsgID(1), []byte(fmt.Sprintf("我再说 %d", i))))
			if err != nil {
				panic(err)
			}
			fmt.Println(c.Write(data))
			data2, err := pack.DefaultPack(messagerr.DefaultMessager(messagerr.MsgID(2), []byte(fmt.Sprintf("我说你 %d", i))))
			if err != nil {
				panic(err)
			}
			fmt.Println(c.Write(data2))
			data2, err = pack.DefaultPack(messagerr.DefaultMessager(messagerr.MsgID(3), []byte(fmt.Sprintf("我说你是个大笨蛋 %d", i))))
			if err != nil {
				panic(err)
			}
			fmt.Println(c.Write(data2))
			fmt.Println(writer.WriteValue([]byte(fmt.Sprintf("我在测第4个 %d  WriteValue", i))))
			fmt.Println(writer.Write(messagerr.MsgID(4), []byte(fmt.Sprintf("我在测第4个 %d  Write", i))))
			fmt.Println(writer.WriteMessager(messagerr.DefaultMessager(messagerr.MsgID(4), []byte(fmt.Sprintf("我在测第4个 %d  WriteMessager", i)))))
			fmt.Println(writer.WriteValue(&msgType{UserName: "axing"}))
			fmt.Println(writer.WriteValue(&ttt.Message{ID: 26}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 27}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 28}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 29}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 30}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 31}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 32}))
		}
	}()

	unmarshaler := encoding.NewUnmarshaler(encoding.Unmarshal(func(mi messagerr.MsgID, b []byte) (interface{}, error) {
		return b, nil
	}), map[messagerr.MsgID]encoding.Unmarshal{5: func(mi messagerr.MsgID, b []byte) (interface{}, error) {
		tmp := &msgType{}
		err := json.Unmarshal(b, tmp)
		// fmt.Println("unm", mi, fmt.Sprintf("%#v", tmp))
		return tmp, err
	}, 6: func(mi messagerr.MsgID, b []byte) (interface{}, error) {
		tmp := &ttt.Message{}
		err := json.Unmarshal(b, tmp)
		// fmt.Println("unm", mi, fmt.Sprintf("%#v", tmp))
		return tmp, err
	}, 7: func(mi messagerr.MsgID, b []byte) (interface{}, error) {
		tmp := &ttttt.Message{}
		err := json.Unmarshal(b, tmp)
		// fmt.Println("unm", mi, fmt.Sprintf("%#v", tmp))
		return tmp, err
	}})
	for {
		idex++
		c, err := list.Accept()
		if err != nil {
			panic(err)
		}
		cc := conn.NewConn(idex, c, conn.WithDispatcher(dispatch.DefaultDispatcher(root, unmarshaler)))
		mmm.AddConn(cc)
		cc.Start()
	}
}
