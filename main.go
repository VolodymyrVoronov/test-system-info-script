package main

import (
	"test-system-info-script/utils"
)

func main() {
	systemInfo := utils.SysInfo{}

	systemInfo.GetCPUInfo()
	systemInfo.GetVmStat()
	systemInfo.GetDiscStat()
	systemInfo.GetHostStat()
	systemInfo.GetNetStat()

	utils.PrintSystemInfo(&systemInfo)
	utils.SaveSystemInfoToFile(&systemInfo)
}
