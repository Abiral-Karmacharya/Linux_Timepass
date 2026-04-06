package storage

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/disk"
)

type PartitionInfo struct {
	Device     string   `json:"device"`
	Mountpoint string   `json:"mountpoint"`
	FSType     string   `json:"fstype"`
	Opts       []string `json:"opts"`
}


func PartitionGet() ([]PartitionInfo, error) {
	var (
		partitionList []PartitionInfo
	)
	partitions, err:= disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving the partitions: %w", err)
	}
	for _, value := range partitions {
		info := PartitionInfo {
			Device:     value.Device,
            Mountpoint: value.Mountpoint,
            FSType:     value.Fstype,
            Opts:       value.Opts,
		}
		partitionList = append(partitionList, info)
	}
	return partitionList, nil
}