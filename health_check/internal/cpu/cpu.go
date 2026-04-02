package cpu

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tui "healthcheck/internal/pkg"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
)

type CoreStats struct {
	AverageUsage float64
	PerCoreUsage []float64
	CriticalCore []int
	WarningCore []int
}

type CPUTemp struct {
	AverageTemp float64
	PerSensorTemp []float64
	WarnCore []int
	CriticalCore []int
	SensorName []string
	Margin []float64
}

type CPUFrequency struct {
	LiveSpeeds []float64
	Average float64
	ModelName string
	PhysCore []int32
	LimitSpeeds []float64
}

type Alert struct {
	CriticalUsage bool
	WarnUsage bool
	CriticalTemp bool
	WarnTemp bool
}

const (
	TjunctionMax = 100.00
	PerCoreUsageWarn = 80.00
	PerCoreUsageCritical = 100.00
	TempWarn = 85.00
	TempCritical = 95.00
)

func getPerCoreStats() (CoreStats, error) {
	var (
		bufferWarning []int
		bufferCritical []int
		totalUsage float64
		averageUsage float64
	)
	corePercent, err := cpu.Percent(time.Second, true)
	if err != nil {
		return CoreStats{}, fmt.Errorf("Error while calculating core usage: %w", err)
	}
	for i, usage := range corePercent {
		if usage > PerCoreUsageWarn {
			bufferWarning = append(bufferWarning, i)
		}
		if usage > PerCoreUsageCritical {
			bufferCritical = append(bufferCritical, i)
		}
		totalUsage += usage
	}
	averageUsage = totalUsage / float64(len(corePercent))
	report := CoreStats{
			AverageUsage: averageUsage,
			PerCoreUsage: corePercent,
			CriticalCore: bufferCritical,
			WarningCore: bufferWarning,
	}	
	return report, nil
}

func getThermal() (CPUTemp, error) {
	var (
		bufferName []string
		bufferTemperature []float64
		bufferMargin []float64
		bufferWarn []int
		bufferCritical []int
		total float64
	)
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		return CPUTemp{}, fmt.Errorf("Error while calculating temperature: %w", err)
	}
	for i, sensor := range sensors {
		margin := TjunctionMax - sensor.Temperature
		if sensor.Temperature > TempWarn {
			bufferWarn = append(bufferWarn, i)
		}
		if sensor.Temperature > TempCritical {
			bufferCritical = append(bufferCritical, i)
		}
		bufferName = append(bufferName, sensor.SensorKey)
		bufferTemperature = append(bufferTemperature, sensor.Temperature)
		bufferMargin = append(bufferMargin, margin)
		total+=sensor.Temperature
	}
	cpuTemp := CPUTemp{
		AverageTemp: total / float64(len(bufferTemperature)),
		PerSensorTemp: bufferTemperature,
		SensorName: bufferName,
		WarnCore: bufferWarn,
		CriticalCore: bufferCritical,
		Margin: bufferMargin,
	}
	return cpuTemp, nil
} 

func checkFrequency() (CPUFrequency, error) {
	info, err := cpu.Info()
	if err != nil {
		return CPUFrequency{}, fmt.Errorf("Error while calculating frequency: %w", err)
	}
	var total float64
	bufferCore := make([]int32, len(info))
	bufferLiveSpeeds := make([]float64, len(info))
	bufferLimit := make([]float64, len(info))
	for i, value := range info {
        data, err:= os.ReadFile(fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cpufreq/scaling_cur_freq", i))
		if err != nil {
			log.Fatalf("Error while reading file: %v", err)
		}
        val, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		if err != nil {
			log.Fatalf("Error while reading file: %v", err)
		}
        bufferLiveSpeeds[i] = val / 1000.00
		bufferCore[i] = value.Cores
		bufferLimit[i] = value.Mhz
		total += val / 1000.0
    }
	average := total / float64(len(bufferLiveSpeeds))
	liveSpeeds := CPUFrequency{
		LiveSpeeds: bufferLiveSpeeds,
		Average: average,
		ModelName: info[0].ModelName,
		PhysCore: bufferCore,
		LimitSpeeds: bufferLimit,
	}
	return liveSpeeds, nil
}

func CPUHealth() {
	style := tui.DefaultStyles()
	coreUsage, err := getPerCoreStats()
	if err != nil {
		log.Fatalf(style.Error.Render("Error in getting core usage: %s"), err)
	}
	coreTemp, err := getThermal()
	if err != nil {
		log.Fatalf(style.Error.Render("Error in getting core temperature: %s"), err)
	}
	coreFrequency, err := checkFrequency()
	if err != nil {
		log.Fatalf(style.Error.Render("Error in getting core frequency: %s"), err)
	}
	fmt.Println(style.Title.Render("-------- CPU Health Check Result --------"))
	fmt.Println(style.Info.Render(coreFrequency.ModelName))
	fmt.Println(style.Normal.Render("Average CPU Core Usage: "), coreUsage.AverageUsage)
	fmt.Println(style.Normal.Render("Average CPU Core Temperature: "), coreTemp.AverageTemp)
	fmt.Println(style.Normal.Render("Average CPU core Frequency: "), coreFrequency.Average)

	if coreUsage.AverageUsage >  PerCoreUsageWarn || coreTemp.AverageTemp > TempWarn{
		fmt.Println(style.Critical.Render("(*) Emergency Alert: Please wait, Calling higher auditing power function."))
		CPUSecurity()
	}
	fmt.Println(style.Info.Render("(+) Everything seems normal. Thank you for using my application ദ്ദി(˵ •̀ ᴗ - ˵ ) ✧"))
}

func CPUSecurity() {
	style := tui.DefaultStyles()
	coreUsage, err := getPerCoreStats()
	if err != nil {
		log.Fatalf(style.Error.Render("Error in getting core usage: %s"), err)
	}
	if len(coreUsage.CriticalCore) > 0  {
		fmt.Println(style.Critical.Render("(*) Crtical core given below, mend them with haste."))
		for _, coreId := range coreUsage.CriticalCore {
			fmt.Printf("Core %d, Usage %.2f\n", coreId, coreUsage.PerCoreUsage[coreId])
		}
	}
	if len(coreUsage.WarningCore) > 0 {
		fmt.Println(style.Warning.Render("(+) Warning core given below, investigate the cause quickly"))
		for _, coreId := range coreUsage.WarningCore{
			fmt.Printf("Core %d, Usage %.2f\n", coreId, coreUsage.PerCoreUsage[coreId])
		}
	}
	coreTemp, err := getThermal()
	if err != nil {
		log.Fatalf(style.Critical.Render("Error in getting core temperature: %s"), err)
	}
	if len(coreTemp.CriticalCore) > 0  {
		fmt.Println(style.Critical.Render("Critical core given below, mend them with haste."))
		for _, coreId := range coreTemp.CriticalCore {
			fmt.Printf("Core %d, Usage %s\n", coreId, coreTemp.SensorName[coreId])
		}
	}
	if len(coreTemp.WarnCore) > 0 {
		fmt.Println(style.Warning.Render("Warning core given below, investigate the case quickly"))
		for _, coreId := range coreTemp.WarnCore {
			fmt.Printf("Core %d, Usage %s\n", coreId, coreTemp.SensorName[coreId])
		}
	} 
}	