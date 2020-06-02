package daemon

import (
	"fmt"
	"github.com/ihaiker/gokit/maths"
	"github.com/shirou/gopsutil/process"
	"strconv"
)

func GetProcessInfo(pid int32) (cupPercent float64, useMem uint64, err error) {
	var p *process.Process
	p, err = process.NewProcess(pid)
	if err != nil {
		return
	}

	var memInfo *process.MemoryInfoStat
	memInfo, err = p.MemoryInfo()
	if err != nil {
		return
	}
	useMem = uint64(maths.Divide64(maths.Add64(float64(memInfo.VMS), float64(memInfo.RSS)), 1024.0*1024.0))

	cupPercent, err = p.CPUPercent()
	if err != nil {
		return
	}

	cupPercent, _ = strconv.ParseFloat(fmt.Sprintf("%0.4f", cupPercent), 10)
	return
}
