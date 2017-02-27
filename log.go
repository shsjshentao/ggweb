package ggweb

import (
	"log"
	"os"
	"strconv"
)

const (
	DefaultDirectory = "/log"
	DefaultFileName  = "all"
	DefaultRotation  = 7
	TRACE            = iota + 1
	DEBUG
	INFO
	ERROR
	PANIC
)

var rotateNum int = 1

type GGLog struct {
	GGLogConfig
	logger *log.Logger
}

type GGLogConfig struct {
	Separation bool
	Level      int
	FileSize   int
	Rotation   int
	Directory  string
	FileName   string
}

func NewDefaultGGLog() GGLog {
	return createGGLog(GGLogConfig{})
}

func NewGGLog(ggLogConfig GGLogConfig) GGLog {
	return createGGLog(ggLogConfig)
}

func createGGLog(ggLogConfig GGLogConfig) GGLog {
	var f *os.File
	var err error
	var filePath string
	if ggLogConfig.FileName == "" {
		ggLogConfig.FileName = DefaultFileName
	}
	if ggLogConfig.Directory == "" {
		ggLogConfig.Directory = DefaultDirectory
	}
	if ggLogConfig.Rotation == 0 {
		ggLogConfig.Rotation = DefaultRotation
	}
	if ggLogConfig.Level == 0 {
		ggLogConfig.Level = INFO
	}

	filePath = ggLogConfig.Directory + ggLogConfig.FileName + "." + strconv.Itoa(rotateNum) + ".log"

	if checkFileIsExist(filePath) {
		f, err = os.OpenFile(filePath, os.O_APPEND, 0666)
	} else {
		f, err = os.Create(filePath)
	}
	if err != nil {
		panic(err)
	}
	var logger = log.New(f, "", log.Ldate|log.Ltime|log.Llongfile)
	ggLog := GGLog{ggLogConfig, logger}
	return ggLog
}

func (g *GGLog) SetFlags(flag int) {
	g.logger.SetFlags(flag)
}

func (g *GGLog) TRACE(i interface{}) {
	g.logger.SetPrefix("[TRACE]")
	g.Output(i, TRACE)
}

func (g *GGLog) DEBUG(i interface{}) {
	g.logger.SetPrefix("[DEBUG]")
	g.Output(i, DEBUG)
}

func (g *GGLog) INFO(i interface{}) {
	g.logger.SetPrefix("[INFO]")
	g.Output(i, INFO)
}

func (g *GGLog) ERROR(i interface{}) {
	g.logger.SetPrefix("[ERROR]")
	g.Output(i, ERROR)
}

func (g *GGLog) PANIC(i interface{}) {
	g.logger.SetPrefix("[PANIC]")
	g.Output(i, PANIC)
}

func (g *GGLog) Output(i interface{}, level int) {
	if level >= g.GGLogConfig.Level {
		g.logger.Println(i)
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
