# ArgsParser

ArgsParser 是使用Golang开发的命令行参数工具，可通过反射自动注册命令行参数。并将
命令行参数注入到结构体（须实现`Arguments`接口）。命令行解析部分使用Golang标准库的
`flag`包实现。


## 使用

#### 安装依赖

```bash
$ go get github.com/irealing/argsparser
```

#### 解析参数

```golang
// 实现Arguments接口
type AppConfig struct{
    X int `param:"x" usage:"int value"`
    Y string `param:"y" usage:"string value"`
}

func (ac *AppConfig) Validate()error{
    if ac.X<0 || ac.Y==""{
        return errors.New("error")
    }
    return nil
}

func main(){
    ac:=&AppConfig{}
    ap:=argsparser.New(ac)

    if err:=ap.Init(); err!=nil{
        fmt.Fatal(err)
    }
    if err:=ac.Parse();err!=nil{
        fmt.Fatal(err)
    }
    fmt.Printf("ac.X: %d\n;ac.Y: %s",ac.X,ac.Y)
}
```
#### 使用CMDHolder

```golang
func main() {
	cfg := new(AppConfig)
	h := argsparser.NewHolder("about")
	h.Register("print", "输出", cfg, func() {
		fmt.Println(cfg)
	})
	h.Execute()
}
```