package main

import (
	"easysalt/command"
	"easysalt/console"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// 服务器配置信息
var configFile = flag.String("c", "", "-c=filename")
// 监听文件地址
var cmd = flag.String("cmd", "tail -f /home/work/logs/nginx/m.access.log", "-cmd=command")
// ssh password
var sshPassword = flag.String("pwd", "", "-pwd=ssh password")
// 配置
var config command.Config

// 初始化
func init() {
	// 解析参数
	flag.Parse()

	// 初始化配置
	servers, err := getServers(*configFile)
	if err != nil {
		fmt.Println(console.ColorfulText(console.TextRed, err.Error()))
		os.Exit(-1)
	}

	config.Cmd = *cmd
	config.Servers = servers
}

func main() {
	var wg sync.WaitGroup
	outputs := make(chan command.Message, 255)

	for _, server := range config.Servers {
		wg.Add(1)
		go func(server command.Server) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf(console.ColorfulText(console.TextRed, "Error: %s\n"), err)
				}
			}()
			defer wg.Done()

			// 填充密码
			if server.Password == "" {
				server.Password = *sshPassword
			}

			cmd := command.NewCommand(server, config.Cmd)
			cmd.Execute(outputs)
		}(server)
	}

	if len(config.Servers) > 0 {
		go func() {
			for output := range outputs {
				content := strings.Trim(output.Content, "\r\n")
				// 去掉文件名称输出
				if content == "" || (strings.HasPrefix(content, "==>") && strings.HasSuffix(content, "<==")) {
					continue
				}

				fmt.Printf(
					"%s %s %s\n",
					console.ColorfulText(console.TextGreen, output.Host),
					console.ColorfulText(console.TextYellow, "->"),
					content,
				)
			}
		}()
	} else {
		fmt.Println(console.ColorfulText(console.TextRed, "No target host is available"))
	}

	wg.Wait()
}

// 获取服务器信息
func getServers(configFile string) ([]command.Server, error) {
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.New("配置文件错误: " + err.Error())
	}

	var servers []command.Server
	err = json.Unmarshal([]byte(config), &servers)
	if err != nil {
		return nil, errors.New("配置解析错误: %v" + err.Error())
	}

	return servers, nil
}
