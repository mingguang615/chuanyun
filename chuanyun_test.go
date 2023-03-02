package chuanyun

import (
	"os"
	"testing"
)

var (
	code   = "AAAAAAAAAAAAAAAAAAAAA"
	secret = "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
	cli    = NewClient(code, secret)
)

//查询业务对象
func TestClient_LoadBizObject(t *testing.T) {
	schemaCode := "D155554ea016ccffaa3426bb20a5193887e3da6"
	bizObjectid := "34b85bb1-e335-489a-a597-e50fa99d1991"
	resp, err := cli.LoadBizObject(schemaCode, bizObjectid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

//批量查询业务对象
func TestClient_LoadBizObjects(t *testing.T) {
	schemaCode := "D155554ea016ccffaa3426bb20a5193887e3da6"
	filter := "{\"FromRowNum\":   0,\"RequireCount\": false,\"ReturnItems\": [],   \"SortByCollection\": [],\"ToRowNum\": 1,   \"Matcher\": { \"Type\": \"And\",   \"Matchers\": []}}"
	resp, err := cli.LoadBizObjects(schemaCode, filter)
	if err != nil {
		t.Fatalf("err:%s", err.Error())
	}
	t.Log(resp)
}

//创建业务对象
func TestClient_CreateBizObject(t *testing.T) {
	schemaCode := "D155554c3d0e079631444918fc8813d21b22d92"
	t.Log(schemaCode)
	body := make(map[string]interface{})
	body["F0000001"] = "测试"
	body["name"] = "测试"
	body["IsSubmit"] = true
	body["ActionName"] = "FenQiDaikuan"
	body["Controller"] = "RongzhiApiController"
	body["AppCode"] = "D155554rongzhi"

	_, err := cli.CreateBizObject(schemaCode, true, &body)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CustomApi(t *testing.T) {
	schemaCode := "D155554c3d0e079631444918fc8813d21b22d92"
	t.Log(schemaCode)
	body := make(map[string]interface{})
	body["name"] = "蜂电"
	body["IsSubmit"] = true
	body["ActionName"] = "FenQiDaikuan"
	body["Controller"] = "RongzhiApiController"
	body["AppCode"] = "D155554rongzhi"

	_, err := cli.CustomApi(&body)
	if err != nil {
		t.Fatal(err)
	}
}

//更新业务对象
func TestClient_UpdateBizObject(t *testing.T) {
	schemaCode := "D155554Fa9bc2edafb4443ef8c29cc7350bd7e83"
	update := map[string]interface{}{
		"F0000007": "testing222",
	}
	resp, err := cli.UpdateBizObject(schemaCode, "40605f54-d98c-468d-ab1a-dae660345b6c", &update)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

//删除业务对象
func TestClient_RemoveBizObject(t *testing.T) {
	schemaCode := "D000886sbtz_"
	_, err := cli.RemoveBizObject(schemaCode, "3cbff021-536f-47fc-a8d4-d30df9448f8f")
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_UploadAnnex(t *testing.T) {
	schemaCode := "D155554ee2ba72b68e646f58102b5fe8676d688"
	fileName := "F0000002"
	bizObjectId := "a8f44b86-124f-48a1-9981-0171ae7005e0"
	file, err := os.Open("./demo.png")
	if err != nil {
		t.Log(err.Error())
		return
	}
	defer file.Close()
	resp, err := cli.UploadAnnex(schemaCode, fileName, bizObjectId, file)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
