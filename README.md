# 氚云OpenApi接口

## Usage

```
go get github.com/mingguang615/chuanyun
```

## Example
```
func main(){
    code   = "AAAAAAAAAAAAAAAAAAAAA"
	secret = "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
	cli    = NewClient(code, secret)
	schemaCode := "D155554ea016ccffaa3426bb20a5193887e3da6"
	bizObjectid := "34b85bb1-e335-489a-a597-e50fa99d1991"
	resp, err := cli.LoadBizObject(schemaCode, bizObjectid)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
```

## API 文档
- 请先阅读官方文档，并获取EngineCode、EngineSecret [氚云OpenAPI官方文档](https://help.h3yun.com/contents/1005/1631.html)
