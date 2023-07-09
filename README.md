![image](https://github.com/Icy2010/Ini-Config/blob/main/ini-config.png)
# Ini-Config
#### 粗略的 Ini文件 操作
#### 作者 Icy 
#### Web http://zelig.cn
```golang

package main

import (
	"fmt"
	z "github.com/Icy2010/Ini-Config"
)

type TContacInfo struct {
	Name   string `ini:"name"`
	Web    string `ini:"web"`
	EMail  string `ini:"email"`
	WeChat string `ini:"wechat"`
	QQ     string `ini:"qq"`
}

func main() {

	type TContacInfo struct {
		Name   string `ini:"name,icy"`
		Web    string `ini:"web,http://zelig.cn"`
		EMail  string `ini:"email,icy2010@hotmail.com"`
		WeChat string `ini:"wechat,IcySoft"`
		QQ     string `ini:"qq,2261206"`
		Title  string `ini:"title,test"`
	}

	ini := NewIniConfig()
	_ = ini.AddSection(`test`).SetStruct(TContacInfo{})
	fmt.Println(ini.GetSection("test"))

	if ini.HasSection(`test`) {
		fmt.Println(`含有 test`)
	}

	ini.ReadFromString(`[default]
string_value = 哈哈 ; 测试1
integer_value = 1 ;测试2 
Float_value = 2.2

[options]
web=https://zelig.cn
name=icy
email=icy2010@hotmail.com
wechat=IcySoft
qq = 2261206
`)

	fmt.Println(ini.GetSection("default"))

	// 读取
	data := TContacInfo{}
	ini.Struct("options", &data)
	fmt.Println(data)

	ini.SaveToFile(`test.ini`)

	ini.ClearSection()
	Sec := ini.AddSection("options")
	Sec.SetString(`web`, `https://zelig.cn`)
	Sec.SetString(`name`, `icy`)
	Sec.SetString(`email`, `icy2010@hotmail.com`)
	Sec.SetString(`wechat`, `IcySoft`)
	Sec.SetString(`qq`, `2261206`)

	fmt.Println(Sec.String("web", ""))

	initext := ""
	ini.SaveToString(&initext)
	fmt.Println("\n" + initext)

	
}

  
```
