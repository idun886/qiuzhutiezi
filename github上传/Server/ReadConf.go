package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

type ServerConf struct {
	Server []struct {
		TunnelPort string `json:"TunnelPort"`
		OpenPort   string `json:"OpenPort"`
	} `json:"Server"`
}

type ConfFileInfo struct {
	FileName    string
	FileContent string
	ConfInfo    *ServerConf
}

func NewConfFileInfo(FileName string) *ConfFileInfo {
	return &ConfFileInfo{
		FileName: FileName,
	}
}

func (c *ConfFileInfo) ReadConfFile() (*ConfFileInfo, error) {
	FileBytes, err := os.ReadFile(c.FileName)
	if err != nil {
		return nil, err
	}
	c.FileContent = string(FileBytes)
	return c, err
}

func (c *ConfFileInfo) ParserConf() (*ServerConf, error) {
	err := json.Unmarshal([]byte(c.FileContent), &c.ConfInfo)
	ConfInfo := c.ConfInfo
	if err != nil {
		errorinfo := "conf parser error :" + err.Error()
		return ConfInfo, errors.New(errorinfo)
	}
	i := 1
	for _, value := range ConfInfo.Server {
		switch true {
		case value.OpenPort == "":
			errorinfo := "Listen moudle " + strconv.Itoa(i) + " error"
			return ConfInfo, errors.New(errorinfo)
		case value.TunnelPort == "":
			errorinfo := "Listen moudle " + strconv.Itoa(i) + " error"
			return ConfInfo, errors.New(errorinfo)
		}
		i++
	}
	return ConfInfo, nil
}
