# feng
本项目使用Go语言实现了一个mini web框架，我将该框架命名为feng。
feng 具备以下五个模块:  
1.用于封装request和response常见方法，在goroutine间传递数据，控制goroutine的Context。  
2.用于HTTP方法匹配，静态路由匹配，动态路由匹配，设置批量通用前缀的路由。  
3.能够在框架中调用中间件的中间件机制。  
4.可以管理模块间关系，降低模块间耦合度的服务和服务容器。  
5.具备应用管理命令和调试模式的命令行工具。  

# 环境配置
开发语言	Golang 1.7.9  
开发工具	Visual Studio Code 1.67.0  
核心类库	net/http  
运行环境	ubuntu20.04  

# 文件目录
/app文件夹存放的是使用框架的一个小型demo。  
/framework文件夹存放的是上述五个功能模块的源码以及对应的测试文件。  
main.go是小型demo的main package。

# 运行方式
1.使用 go build命令编译build_feng_tool文件夹下的feng_tool.go文件，生成命令行工具。  
2.启动命令行工具，该命令工具包括应用管理命令和调试命令。  
3.使用go test命令运行框架的测试用例。  
4.编译main.go文件，运行一个小型demo。  
5.framework文件夹内包含了框架的所有源码，使用人员可自行调用其中内容使用。   








