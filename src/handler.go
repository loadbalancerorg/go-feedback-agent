package main

import (
	"fmt"
	"math"
	"net"
	"strconv"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const (
	Normal = "Normal"
	Halt   = "Halt"
	Drain  = "Drain"
	Down   = "Down"
)

var (
	initialRun = true
)

const (
	returnIdle = true
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	conn.Write(GetResponseForMode())
}

func GetResponseForMode() (response []byte) {
	ramThresholdValue := GlobalConfig.Ram.ThresholdValue.ToFloat()
	cpuThresholdValue := GlobalConfig.Cpu.ThresholdValue.ToFloat()
	cpuImportance := GlobalConfig.Cpu.ImportanceFactor.ToFloat()
	ramImportance := GlobalConfig.Ram.ImportanceFactor.ToFloat()

	switch GlobalConfig.AgentStatus.Value {
	case Normal:
		usedRam := 0.0
		averageCpuLoad := 0.0
		utilization := 0.0
		divider := 0.0

		// Calculate CPU
		if cpuImportance > 0 {
			cpuLoad, err := cpu.Percent(0, false)
			if err != nil {
				return []byte("0%\n")
			}
			averageCpuLoad = cpuLoad[0]
		}

		// Calculate RAM
		if ramImportance > 0 {
			v, err := mem.VirtualMemory()
			if err != nil {
				return []byte("0%\n")
			}
			usedRam = v.UsedPercent
		}

		// If any resource is important and utilized 100% then everything else is not important
		if averageCpuLoad > cpuThresholdValue && cpuThresholdValue > 0 || (usedRam > ramThresholdValue && ramThresholdValue > 0) {
			return []byte("0%\n")
		}

		utilization = utilization + averageCpuLoad*cpuImportance
		if cpuImportance > 0 {
			divider++
		}

		utilization = utilization + usedRam*ramImportance
		if ramImportance > 0 {
			divider++
		}

		for _, tcpService := range GlobalConfig.TCPService {
			// Make sure our importance factor is greater than 0 otherwise ignore
			if tcpService.ImportanceFactor.ToFloat() > 0 {
				// Get session occupied
				sessionOccupied := GetSessionUtilized(tcpService.IPAddress.Value, tcpService.Port.Value, tcpService.MaxConnections.ToInt())

				// Calculate utilization
				utilization = utilization + sessionOccupied*tcpService.ImportanceFactor.ToFloat()

				// increase our divider
				divider++

				if sessionOccupied > 99 && tcpService.ImportanceFactor.ToFloat() == 1 {
					return []byte("0%\n")
				}
			}
		}

		utilization = utilization / divider

		// Account for utilization less than 0
		if utilization < 0 {
			utilization = 0
		}

		// Account for utilization more than 0
		if utilization > 100 {
			utilization = 100
		}

		if returnIdle {
			response = []byte(fmt.Sprintf("%v%%\n", math.Ceil(100-utilization)))
		} else {
			response = []byte(fmt.Sprintf("%v%%\n", math.Ceil(utilization)))
		}

		if initialRun {
			response = append([]byte("up ready "), response...)
		}
	case Drain:
		response = []byte("drain\n")
	case Down:
		response = []byte("down\n")
	case Halt:
		response = []byte("down\n")
	default:
		response = []byte("error\n")
	}

	return
}

func GetSessionUtilized(IPAddress, servicePort string, maxNumberOfSessionsPerService int) (result float64) {
	numberOfEstablishedConnections := getNumberOfLocalEstablishedConnections(IPAddress, servicePort)
	if numberOfEstablishedConnections > 0 && maxNumberOfSessionsPerService > 0 {
		result = float64(numberOfEstablishedConnections) / float64(maxNumberOfSessionsPerService) * 100
	}
	return
}

func getNumberOfLocalEstablishedConnections(ipAddress string, port string) int {
	p, err := strconv.Atoi(port)
	if err != nil {
		return 0
	}

	// get slice of sockets based on match function
	tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		if ipAddress == "*" {
			return s.State == netstat.Established && s.LocalAddr.Port == uint16(p)
		}

		return s.State == netstat.Established && s.LocalAddr.IP.String() == ipAddress && s.LocalAddr.Port == uint16(p)
	})
	if err != nil {
		return 0
	}

	return len(tabs)
}
