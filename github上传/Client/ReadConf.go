package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

type ClientConf struct {
	Client []struct {
		LocalServerAddr string `json:"LocalServerAddr"`
		RemoteAddr      string `json:"RemoteAddr"`
	} `json:"Client"`
}

type ConfFileInfo struct {
	FileName    string
	FileContent string
	ConfInfo    *ClientConf
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

func (c *ConfFileInfo) ParserConf() (*ClientConf, error) {
	err := json.Unmarshal([]byte(c.FileContent), &c.ConfInfo)
	ConfInfo := c.ConfInfo
	if err != nil {
		errorinfo := "conf parser error :" + err.Error()
		return ConfInfo, errors.New(errorinfo)
	}
	i := 1
	for _, value := range ConfInfo.Client {
		switch true {
		case value.RemoteAddr == "":
			errorinfo := "Listen moudle " + strconv.Itoa(i) + " error"
			return ConfInfo, errors.New(errorinfo)
		case value.LocalServerAddr == "":
			errorinfo := "Listen moudle " + strconv.Itoa(i) + " error"
			return ConfInfo, errors.New(errorinfo)
		}
		i++
	}
	return ConfInfo, nil
}
