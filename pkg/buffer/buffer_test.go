package buffer

import (
	"bytes"
	"log"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	rio := bytes.NewBuffer([]byte("987654321"))
	src := []byte("1234567890")
	buf := NewBuffer(rio)
	buf.Write(src)
	args := []struct {
		L      int
		Expect []byte
	}{
		{
			L:      0,
			Expect: []byte{},
		},
		{
			L:      1,
			Expect: []byte("1"),
		},
		{
			L:      2,
			Expect: []byte("23"),
		},
		{
			L:      1,
			Expect: []byte("4"),
		},
		{
			L:      4,
			Expect: []byte("5678"),
		},
		{
			L:      4,
			Expect: []byte("9098"),
		},
		{
			L:      2,
			Expect: []byte("76"),
		},
		{
			L:      5,
			Expect: []byte("54321"),
		},
	}
	for i, v := range args {
		data := make([]byte, v.L)
		n, err := buf.Read(data)
		if err != nil {
			log.Panicln(err)
		}
		if n != v.L || !reflect.DeepEqual(data, v.Expect) {
			log.Printf("TestRead %d read len:%d expect:%#v, n:%d ret:%#v", i, v.L, v.Expect, n, data)
		}
	}
}

func TestReadWrite(t *testing.T) {
	rio := bytes.NewBuffer([]byte("987654321"))
	src := []byte("1234567890")
	buf := NewBuffer(rio)
	buf.Write(src)
	args := []struct {
		L      int
		Expect []byte
	}{
		{
			L:      0,
			Expect: []byte{},
		},
		{
			L:      1,
			Expect: []byte("1"),
		},
		{
			L:      2,
			Expect: []byte("23"),
		},
		{
			L:      1,
			Expect: []byte("4"),
		},
		{
			L:      4,
			Expect: []byte("5678"),
		},
		{
			L:      4,
			Expect: []byte("9098"),
		},
		{
			L:      2,
			Expect: []byte("76"),
		},
		{
			L:      5,
			Expect: []byte("54321"),
		},
	}
	for i, v := range args {
		data := make([]byte, v.L)
		n, err := buf.Read(data)
		if err != nil {
			log.Panicln(err)
		}
		if n != v.L || !reflect.DeepEqual(data, v.Expect) {
			log.Printf("TestRead %d read len:%d expect:%#v, n:%d ret:%#v", i, v.L, v.Expect, n, data)
		}
	}

	// 开始写
	rio.Write(src)
	buf.Write([]byte("ABCDEFG"))

	wargs := []struct {
		L      int
		Expect []byte
	}{
		{
			L:      4,
			Expect: []byte("ABCD"),
		},
		{
			L:      4,
			Expect: []byte("EFG1"),
		},
	}
	for i, v := range wargs {
		data := make([]byte, v.L)
		n, err := buf.Read(data)
		if err != nil {
			log.Panicln(err)
		}
		if n != v.L || !reflect.DeepEqual(data, v.Expect) {
			log.Printf("TestRead %d read len:%d expect:%#v, n:%d ret:%#v", i, v.L, v.Expect, n, data)
		}
	}
}
