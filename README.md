# hot
还有bug, 目标程序出错后不能恢复.
http server 热编译(只能用在windows下)
go build hot.go后, 生成hot.exe, 将hot.exe和hotConf.json复制到你要热编译的目录下. 
修改hotConf.json的内容, filename的值为你mian方法所在文件的文件名. 注意filename不能带后缀名.
修改完毕,双击hot.exe就能自动编译.每次修改go文件后都会重启启动http server.

