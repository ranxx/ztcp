package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/ranxx/ztcp/conn"
	"github.com/ranxx/ztcp/dispatch"
	messagerr "github.com/ranxx/ztcp/message"
	"github.com/ranxx/ztcp/pack"
	"github.com/ranxx/ztcp/router"
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
type message struct {
	id int64
}

func (m message) GetID() int64 {
	return m.id
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
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.New(t)
	vv, ok := v.Interface().(messager)
	fmt.Println(ok)
	vv.GetID()

	fmt.Println(v.MethodByName("GetID").Call(nil)[0].Int())
}

func Test03(t *testing.T) {
	msg := message{id: 4}
	test03(&msg)
	test03(msg)
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
			data, err := pack.NewPackage(messagerr.MsgID(1), []byte(fmt.Sprintf("我再说 %d", i))).PackBytes()
			if err != nil {
				panic(err)
			}
			fmt.Println(c.Write(data))
			data2, err := pack.NewPackage(messagerr.MsgID(2), []byte(fmt.Sprintf("我说你 %d", i))).PackBytes()
			if err != nil {
				panic(err)
			}
			fmt.Println(c.Write(data2))
			data2, err = pack.NewPackage(messagerr.MsgID(3), []byte(fmt.Sprintf("我说你是个大笨蛋 %d", i))).PackBytes()
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
