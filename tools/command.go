package tools

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"publish/sshw"
	"publish/websocket"
)

type Command struct {
	Host string
	Port int
}

func (c *Command) Con() sshw.Client {
	node := new(sshw.Node)
	client := sshw.NewClient(node)
	if c.Host == "" || c.Port <= 0 {
		panic("host or port is null..")
	}
	node.Host = c.Host
	node.Port = c.Port
	return client
}

func (c *Command) RemoteCommand(cmd string) error {
	//当存在其他端口时分割host得到端口
	client := c.Con()
	log.Println("run remote command:", cmd)
	websocket.Broadcast(websocket.Conn, fmt.Sprintf("run remote command:%s\r\n", cmd))
	output, err := client.Cmd(cmd)
	log.Println("\t", string(output))
	websocket.Broadcast(websocket.Conn, fmt.Sprintf("%s\r\n", output))
	return err
}

func (c *Command) LocalCommand(cmd string) error {
	var output []byte
	var err error
	var handel *exec.Cmd
	// 执行单个shell命令时, 直接运行即可
	log.Println("run local command:", cmd)
	websocket.Broadcast(websocket.Conn, fmt.Sprintf("run local command:%s\r\n", cmd))

	handel = exec.Command("/bin/sh", "-c", cmd)
	output, err = handel.Output()
	log.Println("\t", string(output))
	websocket.Broadcast(websocket.Conn, fmt.Sprintf("%s\r\n", output))

	return err
}

func (c *Command) RemoteCommandOutput(cmd string) (output string, err error) {
	//当存在其他端口时分割host得到端口
	client := c.Con()
	log.Println("run remote command:", cmd)
	websocket.Broadcast(websocket.Conn, fmt.Sprintf("run remote command:%s\r\n", cmd))
	result, err := client.Cmd(cmd)
	output = string(result)
	websocket.Broadcast(websocket.Conn, fmt.Sprintf("%s\r\n", output))
	return
}

func (c *Command) LocalCommandOutput(cmd string) (output []byte, err error) {

	var handel *exec.Cmd
	// 执行单个shell命令时, 直接运行即可
	log.Println("run local command:", cmd)
	websocket.Broadcast(websocket.Conn, fmt.Sprintf("run local command:%s\r\n", output))

	handel = exec.Command("/bin/sh", "-c", cmd)
	output, err = handel.Output()
	log.Println("\t", string(output))
	websocket.Broadcast(websocket.Conn, fmt.Sprintf("%s\r\n", output))
	return output, err
}

func (c *Command) FileUpload(sourceFile, destFile string) error {
	client := c.Con()
	output, errput, err := client.UploadFile(sourceFile, destFile)
	log.Println("upload file:", string(output), string(errput))
	websocket.Broadcast(websocket.Conn, fmt.Sprintf("begin upload file: %s\r\n", output))

	return err
}

// 判断文件夹是否存在,不存在则创建
func (c *Command) PathGen(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		log.Println("目录存在：", path)
		return true, nil
	}
	if os.IsNotExist(err) {
		log.Println("目录不存在：", path)
		log.Println("开始创建目录", path)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err.Error())
			return false, err
		}
		return true, err
	}
	return true, err
}

// 判断目录是否存在

func (c *Command) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}
