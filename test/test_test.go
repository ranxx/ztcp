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
	"github.com/ranxx/ztcp/conner"
	"github.com/ranxx/ztcp/dispatch"
	"github.com/ranxx/ztcp/handle"
	"github.com/ranxx/ztcp/pkg/io/write"
	messagerr "github.com/ranxx/ztcp/pkg/message"
	"github.com/ranxx/ztcp/pkg/pack"
	"github.com/ranxx/ztcp/request"
	"github.com/ranxx/ztcp/router"
	"github.com/ranxx/ztcp/server"
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

func test04(i interface{}) {
	fmt.Println(reflect.TypeOf(i).Kind().String())
}

func Test03(t *testing.T) {
	msg := Message{ID: 4}
	test04(&msg)
	test04(msg)
	fmt.Println("-------------------------")
	test04(1)
	i := 10
	test04(i)
	test04(&i)
	// msg2 := ttttt.Message{ID: 5}
	// test03(&msg2)
	// test03(msg2)
	// fmt.Println("-------------------------")
	// msg3 := ttt.Message{ID: 5}
	// test03(&msg3)
	// test03(msg3)
}

func Test04(t *testing.T) {
	mmm := conner.NewManager()
	route1 := router.NewRouter(messagerr.MsgID(1), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("Hhhhhhh", string(s), ok)
	}))
	route2 := router.NewRouter(messagerr.MsgID(2), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("msg-2", string(s), ok)
	}))
	route3 := router.NewRouter(messagerr.MsgID(3), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("msg-3", r.C.RemoteAddr(), string(s), ok)
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
		cc := conn.NewConn(idex, c, conn.WithDispatcher(dispatch.DefaultDispatcher(root)))
		mmm.AddConn(cc)
		cc.Start()
	}
}

func Test05(t *testing.T) {
	mmm := conner.NewManager()
	route1 := router.NewRouter(messagerr.MsgID(1), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("Hhhhhhh", string(s), ok)
	}))
	route2 := router.NewRouter(messagerr.MsgID(2), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("msg-2", string(s), ok)
	}))
	route3 := router.NewRouter(messagerr.MsgID(3), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("msg-3", r.C.RemoteAddr(), string(s), ok)
	}))
	route4 := router.NewRouter(messagerr.MsgID(4), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("msg-4", r.C.RemoteAddr(), string(s), ok)
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
			fmt.Println(writer.WriteBytes(messagerr.MsgID(4), []byte(fmt.Sprintf("我在测第4个 %d  Write", i))))
			fmt.Println(writer.WriteMessager(messagerr.DefaultMessager(messagerr.MsgID(4), []byte(fmt.Sprintf("我在测第4个 %d  WriteMessager", i)))))
		}
	}()
	for {
		idex++
		c, err := list.Accept()
		if err != nil {
			panic(err)
		}
		cc := conn.NewConn(idex, c, conn.WithDispatcher(dispatch.DefaultDispatcher(root)))
		mmm.AddConn(cc)
		cc.Start()
	}
}

func TestStruct(t *testing.T) {
	mmm := conner.NewManager()
	type msgType struct {
		UserName string `json:"user_name"`
	}
	route1 := router.NewRouter(messagerr.MsgID(1), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("Hhhhhhh", string(s), ok)
	}))
	route2 := router.NewRouter(messagerr.MsgID(2), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("msg-2", string(s), ok)
	}))
	route3 := router.NewRouter(messagerr.MsgID(3), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("msg-3", r.C.RemoteAddr(), string(s), ok)
	}))
	route4 := router.NewRouter(messagerr.MsgID(4), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("msg-4", r.C.RemoteAddr(), string(s), ok)
	}))
	route5 := router.NewRouter(messagerr.MsgID(5), handle.WrapHandler(func(c context.Context, r *request.Request) {
		tmp := msgType{}
		err := json.Unmarshal(r.M.GetData(), &tmp)
		log.Println("msg-5", r.C.RemoteAddr(), fmt.Sprintf("%#v", tmp), err)
	}))
	route6 := router.NewRouter(messagerr.MsgID(6), handle.WrapHandler(func(c context.Context, r *request.Request) {
		tmp := ttt.Message{}
		err := json.Unmarshal(r.M.GetData(), &tmp)
		log.Println("msg-6", r.C.RemoteAddr(), fmt.Sprintf("%#v", tmp), err)
	}))
	route7 := router.NewRouter(messagerr.MsgID(7), handle.WrapHandler(func(c context.Context, r *request.Request) {
		tmp := ttttt.Message{}
		err := json.Unmarshal(r.M.GetData(), &tmp)
		log.Println("msg-7", r.C.RemoteAddr(), fmt.Sprintf("%#v", tmp), err)
	}))
	route8 := router.NewRouter(messagerr.MsgID(8), handle.WrapHandler(func(c context.Context, r *request.Request) {
		tmp := ttttt.Message{}
		err := json.Unmarshal(r.M.GetData(), &tmp)
		log.Println("msg-8", r.C.RemoteAddr(), fmt.Sprintf("%#v", tmp), err)
	}))
	route9 := router.NewRouter(messagerr.MsgID(9), handle.WrapHandler(func(c context.Context, r *request.Request) {
		tmp := ttttt.Message{}
		err := json.Unmarshal(r.M.GetData(), &tmp)
		log.Println("msg-9", r.C.RemoteAddr(), fmt.Sprintf("%#v", tmp), err)
	}))
	root := router.NewRoot().AddRouter(route1, route2, route3, route4, route5, route6, route7, route8, route9)

	root.Use(handle.WrapHandler(func(c context.Context, r *request.Request) {
		log.Println("我是mid - 1", r.M.GetMsgID())
	}))
	root.Use(handle.WrapHandler(func(c context.Context, r *request.Request) {
		log.Println("我是mid - 2", r.M.GetMsgID())
		//r.Abort()
	}))
	root.Use(handle.WrapHandler(func(c context.Context, r *request.Request) {
		log.Println("我是mid - 3", r.M.GetMsgID())
	}))

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
			fmt.Println(writer.WriteBytes(messagerr.MsgID(4), []byte(fmt.Sprintf("我在测第4个 %d  Write", i))))
			fmt.Println(writer.WriteMessager(messagerr.DefaultMessager(messagerr.MsgID(4), []byte(fmt.Sprintf("我在测第4个 %d  WriteMessager", i)))))
			fmt.Println(writer.WriteValue(&msgType{UserName: "axing"}))
			fmt.Println(writer.WriteValue(&ttt.Message{ID: 26}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 27}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 28}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 29}))
			fmt.Println(writer.WriteValueWithID(8, &ttttt.Message{ID: 30}))
			fmt.Println(writer.WriteValueWithID(8, &ttttt.Message{ID: 31}))
			fmt.Println(writer.WriteValueWithID(8, &ttttt.Message{ID: 32}))
			fmt.Println(writer.WriteValueWithID(9, &ttttt.Message{ID: 32}))
			fmt.Println(writer.WriteValueWithID(8, &ttttt.Message{ID: 32}))
			fmt.Println(writer.WriteValueWithID(8, &ttttt.Message{ID: 32}))
		}
	}()

	for {
		idex++
		c, err := list.Accept()
		if err != nil {
			panic(err)
		}
		cc := conn.NewConn(idex, c, conn.WithDispatcher(dispatch.DefaultDispatcher(root)))
		mmm.AddConn(cc)
		cc.Start()
	}
}

func TestAbort(t *testing.T) {
	mmm := conner.NewManager()
	type msgType struct {
		UserName string `json:"user_name"`
	}
	route1 := router.NewRouter(messagerr.MsgID(1), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("Hhhhhhh", string(s), ok)
	}))
	route2 := router.NewRouter(messagerr.MsgID(2), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("msg-2", string(s), ok)
	}))
	route3 := router.NewRouter(messagerr.MsgID(3), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("msg-3", r.C.RemoteAddr(), string(s), ok)
	}))
	route4 := router.NewRouter(messagerr.MsgID(4), handle.WrapHandler(func(c context.Context, r *request.Request) {
		s, ok := r.M.GetData(), true
		log.Println("msg-4", r.C.RemoteAddr(), string(s), ok)
	}))
	route5 := router.NewRouter(messagerr.MsgID(5), handle.WrapHandler(func(c context.Context, r *request.Request) {
		tmp := msgType{}
		err := json.Unmarshal(r.M.GetData(), &tmp)
		log.Println("msg-5", r.C.RemoteAddr(), fmt.Sprintf("%#v", tmp), err)
	}))
	route6 := router.NewRouter(messagerr.MsgID(6), handle.WrapHandler(func(c context.Context, r *request.Request) {
		tmp := ttt.Message{}
		err := json.Unmarshal(r.M.GetData(), &tmp)
		log.Println("msg-6", r.C.RemoteAddr(), fmt.Sprintf("%#v", tmp), err)
	}))
	route7 := router.NewRouter(messagerr.MsgID(7), handle.WrapHandler(func(c context.Context, r *request.Request) {
		tmp := ttttt.Message{}
		err := json.Unmarshal(r.M.GetData(), &tmp)
		log.Println("msg-7", r.C.RemoteAddr(), fmt.Sprintf("%#v", tmp), err)
	}))
	route8 := router.NewRouter(messagerr.MsgID(8), handle.WrapHandler(func(c context.Context, r *request.Request) {
		tmp := ttttt.Message{}
		err := json.Unmarshal(r.M.GetData(), &tmp)
		log.Println("msg-8", r.C.RemoteAddr(), fmt.Sprintf("%#v", tmp), err)
	}))
	route9 := router.NewRouter(messagerr.MsgID(9), handle.WrapHandler(func(c context.Context, r *request.Request) {
		tmp := ttttt.Message{}
		err := json.Unmarshal(r.M.GetData(), &tmp)
		log.Println("msg-9", r.C.RemoteAddr(), fmt.Sprintf("%#v", tmp), err)
	}))
	root := router.NewRoot().AddRouter(route1, route2, route3, route4, route5, route6, route7, route8, route9)

	root.Use(handle.WrapHandler(func(c context.Context, r *request.Request) {
		log.Println("我是mid - 1", r.M.GetMsgID())
	}))
	root.Use(handle.WrapHandler(func(c context.Context, r *request.Request) {
		log.Println("我是mid - 2", r.M.GetMsgID())
		// r.Abort()
	}))
	root.Use(handle.WrapHandler(func(c context.Context, r *request.Request) {
		log.Println("我是mid - 3", r.M.GetMsgID())
	}))

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
			fmt.Println(writer.WriteBytes(messagerr.MsgID(4), []byte(fmt.Sprintf("我在测第4个 %d  Write", i))))
			fmt.Println(writer.WriteMessager(messagerr.DefaultMessager(messagerr.MsgID(4), []byte(fmt.Sprintf("我在测第4个 %d  WriteMessager", i)))))
			fmt.Println(writer.WriteValue(&msgType{UserName: "axing"}))
			fmt.Println(writer.WriteValue(&ttt.Message{ID: 26}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 27}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 28}))
			fmt.Println(writer.WriteValue(&ttttt.Message{ID: 29}))
			fmt.Println(writer.WriteValueWithID(8, &ttttt.Message{ID: 30}))
			fmt.Println(writer.WriteValueWithID(8, &ttttt.Message{ID: 31}))
			fmt.Println(writer.WriteValueWithID(8, &ttttt.Message{ID: 32}))
			fmt.Println(writer.WriteValueWithID(9, &ttttt.Message{ID: 32}))
			fmt.Println(writer.WriteValueWithID(8, &ttttt.Message{ID: 32}))
			fmt.Println(writer.WriteValueWithID(8, &ttttt.Message{ID: 32}))
		}
	}()

	for {
		idex++
		c, err := list.Accept()
		if err != nil {
			panic(err)
		}
		cc := conn.NewConn(idex, c, conn.WithDispatcher(dispatch.DefaultDispatcher(root)))
		mmm.AddConn(cc)
		cc.Start()
	}
}

func TestServer(t *testing.T) {
	route := router.NewRouter(0, handle.WrapHandler(func(c context.Context, r *request.Request) {
		log.Println(r.M.GetMsgID(), fmt.Sprintf("%s", r.M.GetData()))
	}))
	route2 := router.NewRouter(2, handle.WrapHandler(func(c context.Context, r *request.Request) {
		log.Println(r.M.GetMsgID(), fmt.Sprintf("%s", r.M.GetData()))
	}))
	root := router.NewRoot().AddRouter(route, route2).NotFound(handle.WrapHandler(func(c context.Context, r *request.Request) {
		log.Println("未知消息", r.M.GetMsgID(), fmt.Sprintf("%s", r.M.GetData()))
	}))

	srv := server.NewServer("tcp", ":12351", server.WithConnOptions(conn.WithDispatcher(dispatch.DefaultDispatcher(root))))

	go func() {
		time.Sleep(time.Second)
		c, err := net.Dial("tcp", ":12351")
		if err != nil {
			panic(err)
		}
		cc := conn.NewConn(0, c)
		writer := cc.Writer()
		i := 0
		for ; ; i++ {
			time.Sleep(time.Second)
			writer.WriteBytes(0, []byte(fmt.Sprintf("我再说 %d", i)))
			writer.WriteBytes(1, []byte(fmt.Sprintf("我再说 %d", i)))
			writer.WriteBytes(2, []byte(fmt.Sprintf("我再说 %d", i)))
			writer.WriteBytes(3, []byte(fmt.Sprintf("我再说 %d", i)))
		}
	}()

	fmt.Println(srv.Start())
}
