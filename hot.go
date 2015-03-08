package main

import (
	"encoding/json"
	"fmt"
	"github.com/hefju/hot/jutil"
	"io/ioutil"
	"os"
	"strings"
    "log"
)

func main() {
    conifg:=InitConfig()
    fmt.Println("file22:", conifg.Filename)
	//HotCompile(filename)
    HotCompile( conifg.Filename)

	// fmt.Println("end")
}

type Config struct {
	Filename string
}
func InitConfig()Config{
    r, err := os.Open("hotConf.json")
    if err != nil {
        log.Fatalln(err)
    }
    decoder := json.NewDecoder(r)
    var c Config
    err = decoder.Decode(&c)
    if err != nil {
        log.Fatalln(err)
    }
    return  c
}

//这个方法会跟txt的编码有关系, 用utf8编码文件名前面会多出一些字符.
func getfilename() string {
	b, err := ioutil.ReadFile("hotConf.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(b))
}
func HotCompile(filename string) {

	// gofile:=filename+".go"
	//
	//    _, err := os.Open(gofile)//"filename.txt"
	//    if err != nil {
	//        fmt.Println("boot failed!\nerror:no file:",gofile, err)
	//        return
	//    }
	//    fmt.Println(filename)
	//    return

	//重启服务器的通道
	restartChan := make(chan int)

	go jutil.SetWatcher(".", restartChan)

	fmt.Println("runner begin")
	runner := jutil.Runner{Filename: filename} //"myhttp"} filename"GoDream"
	runner.Run()
	runner.WaitForRestart(restartChan)
}
