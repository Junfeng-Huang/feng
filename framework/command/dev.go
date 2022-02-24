package command

import (
	"feng/framework/cobra"
	"feng/framework/util"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

func initDevCommand() *cobra.Command {
	return devCommand
}

var devCommand = &cobra.Command{
	Use:   "dev",
	Short: "调试模式",
	Long:  "调试模式，可以方便开发",
	Run: func(c *cobra.Command, args []string) {
		fmt.Println("调试模式启动")
		d := dev{}
		if err := d.serverRestart(); err != nil {
			fmt.Println("启动失败,")
			panic(err)
		}
		go d.monitor()
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		<-quit
		if d.pid != 0 {
			if util.CheckProcessExist(d.pid) {
				syscall.Kill(-d.pid, syscall.SIGKILL)
			}
		}
		fmt.Println("调试模式关闭")
		os.Exit(0)
	},
}

type dev struct {
	folder string
	pid    int
}

func (d *dev) getDevFolder() string {
	if d.folder == "" {
		d.folder = util.GetExecDirectory()
		return d.folder
	}
	return d.folder
}

func (d *dev) serverRestart() error {
	if d.pid != 0 {
		if util.CheckProcessExist(d.pid) {
			err := syscall.Kill(-d.pid, syscall.SIGKILL)
			if err != nil {
				return err
			}
		}
		d.pid = 0
	}
	cmd := exec.Command("go", "run", d.getDevFolder())
	// serverLogFile := filepath.Join(d.folder,"log","app.log")
	// logFile,err := os.OpenFile(serverLogFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC,0644)
	// if err != nil{
	// 	return err
	// }
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		return err
	}
	d.pid = cmd.Process.Pid
	return nil
}

func (d *dev) monitor() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()
	appFolder := d.getDevFolder()
	fmt.Println("监控文件夹:", appFolder)
	filepath.Walk(appFolder, func(path string, info fs.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			return nil
		}
		if util.IsHiddenDirectory(path) {
			return nil
		}
		return watcher.Add(path)
	})
	refreshTime := 3
	t := time.NewTimer(time.Duration(refreshTime) * time.Second)
	t.Stop()
	for {
		select {
		case <-t.C:
			// 计时器到了，代表有文件更新
			fmt.Println("检测到文件更新，重启服务开始")
			if err := d.serverRestart(); err != nil {
				fmt.Println("重启错误:", err)
			}
			fmt.Println("重启服务成功")
			t.Stop()
		case _, ok := <-watcher.Events:
			if !ok {
				continue
			}
			// 有文件更新事件，重置计时器
			t.Reset(time.Duration(refreshTime) * time.Second)
		case _, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			// 监听文件错误
			fmt.Println("监听文件夹错误：", err.Error())
			t.Reset(time.Duration(refreshTime) * time.Second)
		}
	}
}
