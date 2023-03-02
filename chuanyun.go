package chuanyun

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	api             = "https://www.h3yun.com/OpenApi/Invoke"
	upload_file_api = "https://www.h3yun.com/OpenApi/UploadAttachment"
)

type Client struct {
	EngineCode   string
	EngineSecret string
}

type H3Request struct {
	ActionName     string   `json:"ActionName"`
	SchemaCode     string   `json:"SchemaCode"`
	Filter         string   `json:"Filter,omitempty"`
	BizObject      string   `json:"BizObject,omitempty"`
	BizObjectArray []string `json:"BizObjectArray,omitempty"`
	IsSubmit       bool     `json:"IsSubmit"`
	BizObjectId    string   `json:"BizObjectId,omitempty"`
}

type H3Response struct {
	Successful   bool                   `json:"Successful"`
	ErrorMessage interface{}            `json:"ErrorMessage"`
	Logined      bool                   `json:"Logined"`
	ReturnData   map[string]interface{} `json:"ReturnData"`
	DataType     int                    `json:"DataType"`
}

func (r *H3Response) GetReturnData(key string) []byte {
	if r.ReturnData == nil {
		return nil
	}
	data, _ := json.Marshal(r.ReturnData[key])
	return data
}

func (r *H3Response) GetReturnDataMap(key string) (map[string]interface{}, error) {
	if r.ReturnData == nil {
		return nil, errors.New("ReturnData is nil")
	}
	data, _ := json.Marshal(r.ReturnData[key])
	result := make(map[string]interface{})
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/*
Name为字段名，Operator为运算符，Value为数值。

Operator运算符：0 =大于，1=大于等于，2=等于，3=小于等于，4=小于，5=不等于，6=在某个范围内，7=不在某个范围内 8= 模糊查询。
*/
type MatcherItem struct {
	Type     string `json:"Type"`
	Name     string `json:"Name"`
	Operator int    `json:"Operator"`
	Value    string `json:"Value"`
}

type Matchers struct {
	Type     string        `json:"Type"`
	Matchers []MatcherItem `json:"Matchers"`
}

type Filter struct {
	FromRowNum       int      `json:"FromRowNum"`       //分页查询，从第几条开始
	ToRowNum         int      `json:"ToRowNum"`         //分页查询，第几条结束
	RequireCount     bool     `json:"RequireCount"`     //查询的总行数
	ReturnItems      []string `json:"ReturnItems"`      //返回的字段，不填返回所有
	SortByCollection []string `json:"SortByCollection"` //排序字段，目前不支持使用，默认置空
	Matcher          Matchers `json:"Matcher"`
}

func NewFilter() *Filter {
	f := &Filter{
		RequireCount:     false,
		ReturnItems:      make([]string, 0),
		SortByCollection: make([]string, 0),
	}
	f.Matcher.Matchers = make([]MatcherItem, 0)
	f.Matcher.Type = "And"
	return f
}

func (f *Filter) ToString() (string, error) {
	data, err := json.Marshal(&f)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func NewClient(code, secret string) *Client {
	return &Client{
		EngineCode:   code,
		EngineSecret: secret,
	}
}

/*
body: 请求body体
resp：返回结果
*/
func (c *Client) PostRequest(body interface{}, callresp interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("EngineCode", c.EngineCode)
	req.Header.Set("EngineSecret", c.EngineSecret)

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	log.Println("respData: ", string(respData))
	if err != nil {
		return err
	}
	return json.Unmarshal(respData, callresp)
}

/*
desc:LoadBizObject 为加载单个数据，请勿使用该接口来循环加载数据，可以使用LoadBizObjects 来批量加载数据。
actionName：调用的方法名
SchemaCode：需要查询的表单编码
BizObjectId：BizObjectId
*/
func (c *Client) LoadBizObject(schemaCode, bizObjectId string) (*H3Response, error) {
	resp := new(H3Response)
	err := c.PostRequest(&H3Request{
		ActionName:  "LoadBizObject",
		SchemaCode:  schemaCode,
		BizObjectId: bizObjectId,
	}, &resp)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

/*
desc:LoadBizObject 为加载单个数据，请勿使用该接口来循环加载数据，可以使用LoadBizObjects 来批量加载数据。
SchemaCode：需要查询的表单编码
Filter：过滤条件。默认返回前500条数据
*/
func (c *Client) LoadBizObjects(schemaCode, filter string) (*H3Response, error) {
	resp := new(H3Response)
	err := c.PostRequest(&H3Request{
		ActionName: "LoadBizObjects",
		SchemaCode: schemaCode,
		Filter:     filter,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

/*
desc:CreateBizObject 创建单个数据
SchemaCode：表单编码
BizObject：BizObject对象的json 字符串
IsSubmit：为true时创建生效数据，false 为草稿数据
*/
func (c *Client) CreateBizObject(schemaCode string, isSubmit bool, bizObject interface{}) (*H3Response, error) {
	object, err := json.Marshal(bizObject)
	if err != nil {
		return nil, err
	}

	resp := new(H3Response)
	err = c.PostRequest(&H3Request{
		ActionName: "CreateBizObject",
		SchemaCode: schemaCode,
		BizObject:  string(object),
		IsSubmit:   isSubmit,
	}, &resp)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

/*
desc:CreateBizObjects 批量创建业务数据
SchemaCode：表单编码
BizObjects：BizObjects数组对象的json 字符串
IsSubmit：为true时创建生效数据，false 为草稿数据
*/
func (c *Client) CreateBizObjects(schemaCode string, isSubmit bool, bizObjects ...interface{}) (*H3Response, error) {
	arr := make([]string, len(bizObjects))
	for k, v := range arr {
		object, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		arr[k] = string(object)
	}

	resp := new(H3Response)
	err := c.PostRequest(&H3Request{
		ActionName:     "CreateBizObjects",
		SchemaCode:     schemaCode,
		BizObjectArray: arr,
		IsSubmit:       isSubmit,
	}, &resp)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

/*
desc:UpdateBizObject 更新数据。
SchemaCode:表单编码
BizObjectId:表单ObjectId值
BizObject:BizObject对象的json 字符串
*/
func (c *Client) UpdateBizObject(schemaCode, bizObjectId string, bizObject interface{}) (*H3Response, error) {
	object, err := json.Marshal(bizObject)
	if err != nil {
		return nil, err
	}

	resp := new(H3Response)
	err = c.PostRequest(&H3Request{
		ActionName:  "UpdateBizObject",
		SchemaCode:  schemaCode,
		BizObjectId: bizObjectId,
		BizObject:   string(object),
	}, &resp)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

/*
desc: RemoveBizObject 删除表单
SchemaCode:表单编码
BizObjectId:表单ObjectId值
*/
func (c *Client) RemoveBizObject(schemaCode, bizObjectId string) (*H3Response, error) {
	resp := new(H3Response)
	err := c.PostRequest(&H3Request{
		ActionName:  "RemoveBizObject",
		SchemaCode:  schemaCode,
		BizObjectId: bizObjectId,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

/*上传附件
https://www.h3yun.com/OpenApi/UploadAttachment?SchemaCode=D000024chuangjian&FilePropertyName=F0000011&BizObjectId=7140d02a-6d19-461b-b7dc-892f202b1566
*/
func (c *Client) UploadAnnex(schemaCode, filePropertyName, bizObjectId string, file *os.File) (*H3Response, error) {
	if file == nil {
		return nil, errors.New("file is not nil")
	}

	api := fmt.Sprintf("%s?SchemaCode=%s&FilePropertyName=%s&BizObjectId=%s",
		upload_file_api, schemaCode, filePropertyName, bizObjectId)

	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)

	// file part1
	_, fileName := filepath.Split(file.Name())
	fw1, err := bw.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	io.Copy(fw1, file)
	bw.Close() //write the tail boundry

	req, err := http.NewRequest("POST", api, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("EngineCode", c.EngineCode)
	req.Header.Set("EngineSecret", c.EngineSecret)
	req.Header.Set("Content-Type", bw.FormDataContentType())

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := new(H3Response)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

//自定义API接口
func (c *Client) CustomApi(bizObject interface{}) (*H3Response, error) {
	resp := new(H3Response)
	err := c.PostRequest(bizObject, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
