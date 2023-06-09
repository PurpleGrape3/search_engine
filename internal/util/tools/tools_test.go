package tools

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type Data struct {
	Modified  string
	Saled     string
	CreatedAt time.Time
	Num       int
}
type Doc struct {
	Ident string
	Data
}

func TestBytesStr(t *testing.T) {
	str := "hello world"
	str = ""
	b := Str2Bytes(str)
	t.Log(b)
	str2 := Bytes2Str(b)
	t.Log(str2)
}

func TestConvStruct2Map(t *testing.T) {
	doc := Doc{Ident: "88.199.1/abc.def", Data: Data{Modified: "北京市丰台区", Saled: "北京市海淀区", CreatedAt: time.Now()}}
	objMap, err := ConvStruct2Map(doc)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(objMap)

	doc = Doc{Ident: "88.199.1/abc.def", Data: Data{Modified: "北京市丰台区"}}
	objMap, err = ConvStruct2Map(doc)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(objMap)
}

func TestInterfaceLt(t *testing.T) {
	var a interface{} = nil
	var b interface{} = 15
	var s interface{} = "hello"
	ret, err := InterfaceLt(a, b)
	t.Log(ret, err)
	ret, err = InterfaceLt(a, s)
	t.Log(ret, err)
}

func TestJson(t *testing.T) {
	m := make(map[string]interface{})
	doc := Doc{Ident: "88.199.1/abc.def", Data: Data{Modified: "北京市丰台区", Saled: "北京市海淀区", CreatedAt: time.Now(), Num: 10}}
	docByte, _ := json.Marshal(doc)
	t.Log(string(docByte))
	err := json.Unmarshal(docByte, &m)
	if err != nil {
		t.Log(err)
	}
	t.Log(m, reflect.TypeOf(m["Num"]), reflect.TypeOf(m["CreatedAt"]))
}
