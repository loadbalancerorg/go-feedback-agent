package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"math"
	"net"
	"strings"
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

func handleClient(conn net.Conn) {
	_, err := conn.Write(GetResponseForMode())
	if err != nil {
		eventLog.Error("Feedback Agent: Failed to write response")
	}
	err = conn.Close()
	if err != nil {
		eventLog.Error("Feedback Agent: Failed to close connection")
	}
}

func GetResponseForMode() []byte {

	response := []byte("")
	switch GlobalConfig.AgentStatus.Value {
	case Normal:
		eventLog.Debug("Feedback Agent: Normal mode")
		utilization, err := CalculateNormalState()
		if err != nil {
			response = []byte("error\n")
		} else {
			if GlobalConfig.ReturnIdle.Value == "true" || GlobalConfig.ReturnIdle.Value == "" {
				response = []byte(fmt.Sprintf("%v%%\n", math.Ceil(100-utilization)))
			} else {
				response = []byte(fmt.Sprintf("%v%%\n", math.Ceil(utilization)))
			}
		}
	case Drain:
		eventLog.Debug("Feedback Agent: Normal drain")
		response = []byte("drain\n")
	case Down:
		eventLog.Debug("Feedback Agent: Normal down")
		response = []byte("down\n")
	case Halt:
		eventLog.Debug("Feedback Agent: Normal halt")
		response = []byte("down\n")
	default:
		eventLog.Debug("Feedback Agent: Normal error")
		response = []byte("error\n")
	}

	if initialRun {
		response = append([]byte("up ready "), response...)
	}

	return response
}

func CalculateNormalState() (float64, error) {

	eventLog.Debug("Feedback Agent: CalculateNormalState")
	cpuImportance := GlobalConfig.Cpu.ImportanceFactor.ToFloat()
	ramImportance := GlobalConfig.Ram.ImportanceFactor.ToFloat()

	averageCpuLoad, usedRam, err := SetImportanceValues(cpuImportance, ramImportance)
	if err != nil {
		eventLog.Debug("Feedback Agent: Error get Importance Values")
		return 0, err
	}
	eventLog.Debug(fmt.Sprintf("Feedback Agent: averageCpuLoad = %f, usedRam = %f", averageCpuLoad, usedRam))

	ramThresholdValue := GlobalConfig.Ram.ThresholdValue.ToFloat()
	cpuThresholdValue := GlobalConfig.Cpu.ThresholdValue.ToFloat()
	eventLog.Debug(fmt.Sprintf("Feedback Agent: ramThresholdValue = %f, cpuThresholdValue = %f", ramThresholdValue, cpuThresholdValue))

	// If any resource is important and utilized 100% then everything else is not important
	if (averageCpuLoad > cpuThresholdValue && cpuThresholdValue > 0) ||
		(usedRam > ramThresholdValue && ramThresholdValue > 0) {
		eventLog.Debug("Feedback Agent: important override")
		return 0, nil
	}

	utilization, err := CalculateUtilization(averageCpuLoad, cpuImportance, usedRam, ramImportance)
	if err != nil {
		eventLog.Debug("Feedback Agent: Error Calculate Utilization")
		return 0, err
	}

	return utilization, nil
}

func CalculateUtilization(averageCpuLoad float64, cpuImportance float64, usedRam float64, ramImportance float64) (float64, error) {
	utilizationSystem := getSystemUtilization(averageCpuLoad, cpuImportance, usedRam, ramImportance)
	eventLog.Debug(fmt.Sprintf("Feedback Agent(CalculateUtilization): utilizationSystem = %f", utilizationSystem))

	utilizationServices := getServicesUtilization()
	eventLog.Debug(fmt.Sprintf("Feedback Agent(CalculateUtilization): utilizationServices = %f", utilizationServices))
	if utilizationServices == 100 {
		eventLog.Debug("Feedback Agent(CalculateUtilization): utilizationServices == 100")
		return 100, nil
	}

	utilization := utilizationSystem + utilizationServices
	eventLog.Debug(fmt.Sprintf("Feedback Agent(CalculateUtilization): utilization = %f", utilization))

	// Account for utilization less than 0
	if utilization < 0 {
		utilization = 0
	}

	// Account for utilization more than 0
	if utilization > 100 {
		utilization = 100
	}
	return utilization, nil
}

func getSystemUtilization(averageCpuLoad float64, cpuImportance float64, usedRam float64, ramImportance float64) float64 {
	divider := 0.0
	utilizationSystem := 0.0
	if cpuImportance > 0 {
		utilizationSystem = utilizationSystem + (averageCpuLoad * cpuImportance)
		divider++
	}
	if cpuImportance == 1 && averageCpuLoad > 99 {
		eventLog.Debug("Feedback Agent(getSystemUtilization): cpuImportance == 1 && averageCpuLoad > 99")
		return 100
	}

	if ramImportance > 0 {
		utilizationSystem = utilizationSystem + (usedRam * ramImportance)
		divider++
	}
	if ramImportance == 1 && usedRam > 99 {
		eventLog.Debug("Feedback Agent(getSystemUtilization): ramImportance == 1 && usedRam > 99")
		return 100
	}

	eventLog.Debug(fmt.Sprintf("Feedback Agent(getSystemUtilization): utilizationSystem = %f, divider = %f", utilizationSystem, divider))

	if divider > 0 {
		utilizationSystem = utilizationSystem / divider
	}
	return utilizationSystem
}

func getServicesUtilization() float64 {
	eventLog.Debug("Feedback Agent(getServicesUtilization): Start Services calculation")

	utilization := 0.0
	divider := 0.0
	for _, tcpService := range GlobalConfig.TCPService {
		eventLog.Debug(fmt.Sprintf("Feedback Agent(getSystemUtilization): TCPService = %s", tcpService.Name))

		if tcpService.ImportanceFactor.ToFloat() == 0 {
			eventLog.Debug("Feedback Agent(getServicesUtilization): Service not important")
			continue
		}
		sessionOccupied := 100.0
		if tcpService.MaxConnections.ToInt() > 0 {
			eventLog.Debug(fmt.Sprintf("Feedback Agent(getSystemUtilization): MaxConnections = %d", tcpService.MaxConnections.ToInt()))
			// Get session occupied
			numberOfEstablishedConnections := getNumberOfLocalEstablishedConnections(tcpService.IPAddress.Value, tcpService.Port.Value)
			eventLog.Debug(fmt.Sprintf("Feedback Agent(getSystemUtilization): numberOfEstablishedConnections = %d", numberOfEstablishedConnections))
			sessionOccupied = float64(numberOfEstablishedConnections) / float64(tcpService.MaxConnections.ToInt()) * 100
			eventLog.Debug(fmt.Sprintf("Feedback Agent(getSystemUtilization): sessionOccupied = %f", sessionOccupied))
		}

		if sessionOccupied > 99 && tcpService.ImportanceFactor.ToFloat() == 1 {
			eventLog.Debug("Feedback Agent(getServicesUtilization): sessionOccupied > 99 && tcpService.ImportanceFactor = 1")
			return 100
		}

		// Calculate utilization
		utilization = utilization + sessionOccupied*tcpService.ImportanceFactor.ToFloat()

		// increase our divider
		divider++

	}

	eventLog.Debug(fmt.Sprintf("Feedback Agent(getSystemUtilization): utilization = %f, divider = %f", utilization, divider))

	if divider > 0 {
		utilization = utilization / divider
	}
	eventLog.Debug(fmt.Sprintf("Feedback Agent(getSystemUtilization): utilization / divider = %f", utilization))

	return utilization
}

func SetImportanceValues(cpuImportance float64, ramImportance float64) (float64, float64, error) {
	// Calculate CPU
	averageCpuLoad := 0.0
	if cpuImportance > 0 {
		cpuLoad, err := cpu.Percent(0, false)
		if err != nil {
			return 0, 0, err
		}
		averageCpuLoad = cpuLoad[0]
	}

	// Calculate RAM
	usedRam := 0.0
	if ramImportance > 0 {
		v, err := mem.VirtualMemory()
		if err != nil {
			return 0, 0, err
		}
		usedRam = v.UsedPercent
	}
	eventLog.Debug(fmt.Sprintf("Feedback Agent(SetImportanceValues): averageCpuLoad = %f, usedRam = %f", averageCpuLoad, usedRam))

	return averageCpuLoad, usedRam, nil
}

func getNumberOfLocalEstablishedConnections(ipAddress string, port string) int {
	if ipAddress == "*" {
		ipAddress = ""
	}
	result := runcmd("netstat -nt | findstr " + ipAddress + ":" + port + "  | findstr ESTABLISHED ")
	count := len(strings.Split(result, "\n"))

	eventLog.Debug(fmt.Sprintf("Feedback Agent(getNumberOfLocalEstablishedConnections): count = %d", count))

	if count == 0 {
		return 0
	}
	return count - 1
}
