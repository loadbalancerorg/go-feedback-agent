package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var currentAgentStatus string = ""
var GlobalConfig *XMLConfig

type LogsLevel int32

const (
	LogsINFO    LogsLevel = 0
	LogsWARNING LogsLevel = 1
	LogsERROR   LogsLevel = 2
	LogsDEBUG   LogsLevel = 3
)

type ValueAttr struct {
	Value string `xml:"value,attr"`
}

func (va ValueAttr) ToInt() int {
	val, err := strconv.Atoi(va.Value)
	if err != nil {
		panic(err)
	}
	return val
}

func (va ValueAttr) ToFloat() float64 {
	val, err := strconv.ParseFloat(va.Value, 64)
	if err != nil {
		panic(err)
	}
	return val
}
func (va ValueAttr) ToString() string {
	return va.Value
}

type TCPService struct {
	Name             ValueAttr
	IPAddress        ValueAttr
	Port             ValueAttr
	MaxConnections   ValueAttr
	ImportanceFactor ValueAttr
}

type CPU struct {
	ImportanceFactor ValueAttr
	ThresholdValue   ValueAttr
}

type RAM struct {
	ImportanceFactor ValueAttr
	ThresholdValue   ValueAttr
}

type XMLConfig struct {
	XMLName                           xml.Name `xml:"xml"`
	Cpu                               CPU
	Ram                               RAM
	TCPService                        []TCPService
	ReadAgentStatusFromConfig         ValueAttr
	ReadAgentStatusFromConfigInterval ValueAttr
	AgentStatus                       ValueAttr
	Interval                          ValueAttr
	Port                              ValueAttr
	ReturnIdle                        ValueAttr
	LogLevel                          ValueAttr
}

func readConfig() {
	eventLog.Info("Feedback Agent: Reading Config")
	xmlFile, err := os.Open(CONFIG_FILE)
	if err != nil {
		panic(err)
	}
	defer xmlFile.Close()
	content, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(content, &GlobalConfig)
	if err != nil {
		panic(err)
	}
}

func InitConfig() {
	setupLogging()

	readConfig()

	intervalTicker := time.NewTicker(time.Second * time.Duration(GlobalConfig.Interval.ToInt()))
	go func() {
		for {
			select {
			case <-intervalTicker.C:
				initialRun = false
			}
		}
	}()

	if strings.ToLower(GlobalConfig.ReadAgentStatusFromConfig.Value) == "true" {
		statusTicker := time.NewTicker(time.Second * time.Duration(GlobalConfig.ReadAgentStatusFromConfigInterval.ToInt()))
		go func() {
			for {
				select {
				case <-statusTicker.C:
					readConfig()
					// If status changed, send 'up ready' for a full interval
					if currentAgentStatus != GlobalConfig.AgentStatus.Value {
						initialRun = true
						intervalTicker.Stop()
						intervalTicker = time.NewTicker(time.Second * time.Duration(GlobalConfig.Interval.ToInt()))
						go func() {
							for {
								select {
								case <-intervalTicker.C:
									initialRun = false
								}
							}
						}()
					}
					currentAgentStatus = GlobalConfig.AgentStatus.Value
				}
			}
		}()
	}
}
