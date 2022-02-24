package command

import (
	"feng/framework/cobra"
	"feng/framework/util"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

// app start flags var
var mainFolder = ""

// var appDaemon = false // 是否以daemon方式启动

// 初始化app命令及其子命令
func initAppCommand() *cobra.Command {
	appStartCommand.Flags().StringVarP(&mainFolder, "main", "m", "", "设置package main路径")
	// appStartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", false, "以daemon方式启动app")

	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appStateCommand)
	appCommand.AddCommand(appRestartCommand)
	appCommand.AddCommand(appStopCommand)
	return appCommand
}

var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，包含业务启动，关闭，重启，查询等功能",
	RunE: func(c *cobra.Command, args []string) error {
		// 打印帮助文档
		c.Help()
		return nil
	},
}

// appStartCommand 启动一个Web服务
var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个app服务",
	RunE: func(c *cobra.Command, args []string) error {
		if mainFolder == "" {
			mainFolder = util.GetExecDirectory()
		}
		pidFolder := filepath.Join(util.GetExecDirectory(), "runtime")
		if !util.Exists(pidFolder) {
			if err := os.Mkdir(pidFolder, os.ModePerm); err != nil {
				return err
			}
		}
		serverPidFile := filepath.Join(pidFolder, "app.pid")
		logFolder := filepath.Join(util.GetExecDirectory(), "log")
		if !util.Exists(logFolder) {
			if err := os.Mkdir(logFolder, os.ModePerm); err != nil {
				return err
			}
		}
		serverLogFile := filepath.Join(logFolder, "app.log")

		// if appDaemon {
		// 	// 创建一个Context
		// 	cntxt := &daemon.Context{
		// 		// 设置pid文件
		// 		PidFileName: serverPidFile,
		// 		PidFilePerm: 0664,
		// 		// 设置日志文件
		// 		LogFileName: serverLogFile,
		// 		LogFilePerm: 0640,
		// 		// 设置工作路径
		// 		WorkDir: mainFolder,
		// 		// 设置所有设置文件的mask，默认为750
		// 		Umask: 027,
		// 		// 子进程的参数，按照这个参数设置，子进程的命令为 ./hade app start --daemon=true
		// 		Args: []string{"", "app", "start", "--daemon=true"},
		// 	}
		// 	// 启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
		// 	d, err := cntxt.Reborn()
		// 	if err != nil {
		// 		return err
		// 	}
		// 	if d != nil {
		// 		// 父进程直接打印启动成功信息，不做任何操作
		// 		fmt.Println("app启动成功，pid:", d.Pid)
		// 		fmt.Println("日志文件:", serverLogFile)
		// 		return nil
		// 	}
		// 	defer cntxt.Release()
		// 	// 子进程执行真正的app启动操作
		// 	fmt.Println("deamon started")

		// 	if err := startServe(mainFolder); err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	return nil
		// }
		// 非daemon模式直接执行
		cmdName, err := exec.LookPath("go")
		if err != nil {
			log.Fatalln("feng go:请先安装Go")
		}
		params := []string{"run"}
		params = append(params, mainFolder)
		logFile, err := os.OpenFile(serverLogFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		cmd := exec.Command(cmdName, params...)
		cmd.Stdout = logFile
		// cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		err = cmd.Start()
		if err != nil {
			fmt.Println("go run err:")
			fmt.Println(err)
			fmt.Println("----------------")
			return err
		}
		fmt.Println("Feng App开始运行")
		pid := cmd.Process.Pid
		fmt.Println("[PID]", pid)
		err = ioutil.WriteFile(serverPidFile, []byte(strconv.Itoa(pid)), 0644)
		if err != nil {
			return err
		}
		go func() {
			cmd.Wait()
		}()
		return err
	},
}

// 重新启动一个app服务
var appRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重新启动一个app服务",
	RunE: func(c *cobra.Command, args []string) error {
		pidFolder := filepath.Join(util.GetExecDirectory(), "runtime")
		serverPidFile := filepath.Join(pidFolder, "app.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if content != nil && len(content) > 1 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				if err := syscall.Kill(-pid, syscall.SIGKILL); err != nil {
					return err
				}
				for {
					if !util.CheckProcessExist(pid) {
						break
					}
					fmt.Println("正在关闭进程")
					time.Sleep(2 * time.Second)
				}
				if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
					return err
				}
				fmt.Println("成功停止app服务进程[PID]:", pid)

			}
		}
		return appStartCommand.RunE(c, args)

	},
}

// 停止一个已经启动的app服务
var appStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止一个已经启动的app服务",
	RunE: func(c *cobra.Command, args []string) error {
		pidFolder := filepath.Join(util.GetExecDirectory(), "runtime")
		serverPidFile := filepath.Join(pidFolder, "app.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if content != nil && len(content) > 1 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				if err := syscall.Kill(-pid, syscall.SIGKILL); err != nil {
					return err
				}
				if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
					return err
				}

				fmt.Println("停止app服务进程[PID]:", pid)
				return nil
			}
		}
		fmt.Println("没有app服务进程存在")
		return nil
	},
}

// 获取启动的app的pid
var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "获取启动的app的pid",
	RunE: func(c *cobra.Command, args []string) error {
		pidFolder := filepath.Join(util.GetExecDirectory(), "runtime")
		serverPidFile := filepath.Join(pidFolder, "app.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if content != nil && len(content) > 1 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				fmt.Println("app服务已经启动, pid:", pid)
				return nil
			}
		}
		fmt.Println("没有app服务存在")
		return nil
	},
}
