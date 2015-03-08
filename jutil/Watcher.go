package jutil
import (
	"fmt"
	"io/ioutil"
	"github.com/go-fsnotify/fsnotify"
	"log"
	"path"
    "os/exec"
    "time"
    "os"
)


var lstFolder []string

func SetWatcher(root string, order chan int) {
    lstFolder = make([]string, 0)//存放监视文件夹的slice
    lstFolder = append(lstFolder, root) //首先就把当前目录加入slice
    GetFoleder(root)//递归把所有子文件夹都加入到数组中

    watcher, err := fsnotify.NewWatcher()//文件改动监视器
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    done := make(chan bool)
    go func() {
        for {
            select {
            case event := <-watcher.Events:     //如果有go文件修改,才触发更新
                if path.Ext(event.Name) == ".go" && event.Op&fsnotify.Write == fsnotify.Write {
                    //fmt.Println(event.Name, "--", event.Op, "--", event)
                    order <- 0  //触发重启事件
                }

            case err := <-watcher.Errors:
                log.Println("watcher.Errors:", err)
            }
        }
    }()
    fmt.Println("watcher on-line.") //开始监视

    for _, f := range lstFolder {
        err = watcher.Add(f)
        if err != nil {
            log.Fatal("watcher.Add:", err)
        }
    }

    <-done
}

//遍历文件夹, 遇到以.开头的文件夹就跳过, 不知道如何检查隐藏文件夹,如果能的话都跳过
func GetFoleder(root string) {
    files, _ := ioutil.ReadDir(root)
    for _, fi := range files {
        if fi.IsDir() {
            if fi.Name()[0] == '.' { //不处理以.开头的文件夹
                continue
            }
            mypath := path.Join(root, fi.Name())
            lstFolder = append(lstFolder, mypath)
            GetFoleder(mypath)
        }
    }
}


//用来重启 http server
//负责go run (http server) 和 taskkill正在运行的http server
type Runner struct {
    Filename string  //准备运行的文件名
    lastExec time.Time //上次执行时间, 防止连续执行程序
}

//重启服务器
func (x Runner) WaitForRestart(order chan int) {
    for {
        <-order
        now:=time.Now()//检验时间, 如果超过1秒才重启http server
        if now.After(x.lastExec.Add(time.Second)){
            x.lastExec=time.Now()
            x.Kill()
            x.Run()
           // fmt.Println("restart:" + x.Filename)
        }
    }
}

//执行命令
func (x Runner) Run() {
    goname := x.Filename + ".go"
    fmt.Println("Run:" + goname)
    c := exec.Command("go", "run", goname) // "myserver.go")
    c.Stdout=os.Stdout
    c.Start()
}

//终结进程.
func (x Runner) Kill() {
    taskname := x.Filename + ".exe"
    fmt.Println("Kill:" + taskname)
    c := exec.Command("taskkill.exe", "/f", "/im", taskname) //"myserver.exe")
    err := c.Start()
    if err != nil {
        fmt.Println(err)
    }
}

