package memory

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	tui "healthcheck/internal/pkg"

	"github.com/shirou/gopsutil/v3/mem"
)


type BasicMemory struct {
	Total uint64
	Available uint64
	UsedPercent float64
	SwapTotal uint64
	Status int
}

type ECCStatus struct {
	Present bool
	Errors int	
}

var style = tui.DefaultStyles()
const (
	warnMem = 70
	criticalMem = 90
	StatusOK = 0
	StatusWarning = 1
	StatusCritical = 2
)

func GetMemory() (BasicMemory, error) {
	var (
		status int
		bufferUsage float64
	)
	memory, err := mem.VirtualMemory()
	if err != nil {
		return BasicMemory{}, fmt.Errorf("Error while retrieving basic info: %w", err)
	}
	status = StatusOK
	bufferUsage = memory.UsedPercent
	if bufferUsage > warnMem {
		status = StatusWarning
	}else if bufferUsage > criticalMem{
		status = StatusCritical
	}
	report := BasicMemory{
		Total: memory.Total,
		Available: memory.Available,
		UsedPercent: memory.UsedPercent,
		SwapTotal: memory.SwapTotal,
		Status: status,

	}
	return report, nil
}

func GetECC() ECCStatus {
	data, err := os.ReadFile("/sys/devices/system/edac/mc/mc0/ce_count")
	if err != nil {
		return ECCStatus{Present:false, Errors: 0}
	}
	val, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return ECCStatus{Present:true, Errors:val}
}

func MemoryHealth() {
	var (
		wg sync.WaitGroup
		basicChan = make(chan BasicMemory, 1)
		eccChan = make(chan ECCStatus, 1)
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		res, err := GetMemory()
		if err != nil {
			log.Fatalf(style.Error.Render("%s"), err)
			return 
		}
		basicChan <- res
	}()
	go func() {
		wg.Done()
		res := GetECC()
		eccChan <- res
	}()
	basicMemory := <-basicChan
	eccStatus := <-eccChan
	close(basicChan)
	close(eccChan)
	basicEnum := reflect.ValueOf(basicMemory)
	fmt.Println(style.Title.Render("----- Memory Health Check Report -----"))
	for i := 0; i < basicEnum.NumField(); i++ {
		fmt.Println(style.Normal.Render(basicEnum.Type().Field(i).Name) , basicEnum.Field(i))
	}
	if !eccStatus.Present {
		fmt.Println(style.Normal.Render("(*) Linus torvald would face palm. You don't have ECC"))
	}else {
		fmt.Println(style.Normal.Render("(+) You are one step closer to linus torvalds. You have ECC wow"))
		if eccStatus.Errors > 0 {
			fmt.Println(style.Warning.Render("(*) There are some error with ECC: "),eccStatus.Errors)
		}else {
			fmt.Println(style.Normal.Render("(+) No error in your ECC too"))
		}
	}
	if basicMemory.Status == StatusCritical {
		fmt.Println(style.Critical.Render("(*) Emergency Alert: Please wait, Calling higher auditing power function."))
	}else {
		fmt.Println(style.Info.Render("Your Memory is okay ₍^. .^₎⟆"))
	}
}