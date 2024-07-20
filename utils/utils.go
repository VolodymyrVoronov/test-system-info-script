package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type SysInfo struct {
	CPUInfo  []cpu.InfoStat         `json:"cpu_info"`
	VmStat   *mem.VirtualMemoryStat `json:"mem_info"`
	DiscStat *disk.UsageStat        `json:"disk_info"`
	HostStat *host.InfoStat         `json:"host_info"`
	NetStat  []net.IOCountersStat   `json:"net_info"`
}

func (s *SysInfo) GetCPUInfo() {
	var cpus []cpu.InfoStat

	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Fatalf("Error fetching CPU info: %v", err)
		return
	}

	cpus = append(cpus, cpuInfo...)

	s.CPUInfo = cpus
}

func (s *SysInfo) GetVmStat() {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Error fetching memory info: %v", err)
		return
	}

	s.VmStat = vmStat
}

func (s *SysInfo) GetDiscStat() {
	diskStat, err := disk.Usage("/")
	if err != nil {
		log.Fatalf("Error fetching disk info: %v", err)
		return
	}

	s.DiscStat = diskStat
}

func (s *SysInfo) GetHostStat() {
	hostStat, err := host.Info()
	if err != nil {
		log.Fatalf("Error fetching host info: %v", err)
		return
	}

	s.HostStat = hostStat
}

func (s *SysInfo) GetNetStat() {
	var nets []net.IOCountersStat

	netStat, err := net.IOCounters(true)
	if err != nil {
		log.Fatalf("Error fetching network info: %v", err)
		return
	}

	nets = append(nets, netStat...)

	s.NetStat = nets
}

func PrintSystemInfo(s *SysInfo) {
	color.New(color.FgCyan, color.Bold).Println("\nSystem Information")

	color.New(color.FgGreen).Println("\nHost Information:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Property", "Value"})

	table.Append([]string{"Hostname", s.HostStat.Hostname})
	table.Append([]string{"OS", s.HostStat.OS})
	table.Append([]string{"Platform", s.HostStat.Platform})
	table.Append([]string{"Platform Family", s.HostStat.PlatformFamily})
	table.Append([]string{"Platform Version", s.HostStat.PlatformVersion})
	table.Append([]string{"Kernel Version", s.HostStat.KernelVersion})
	table.Append([]string{"Kernel Arch", s.HostStat.KernelArch})
	table.Append([]string{"Uptime", fmt.Sprintf("%d seconds", s.HostStat.Uptime)})
	table.Append([]string{"Boot Time", fmt.Sprintf("%d", s.HostStat.BootTime)})
	table.Append([]string{"Procs", fmt.Sprintf("%d", s.HostStat.Procs)})
	table.Render()

	color.New(color.FgGreen).Println("\nCPU Information:")
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"CPU", "Model Name", "Cores", "Mhz"})
	for _, cpu := range s.CPUInfo {
		table.Append([]string{
			fmt.Sprintf("%d", cpu.CPU),
			cpu.ModelName,
			fmt.Sprintf("%d", cpu.Cores),
			fmt.Sprintf("%.2f", cpu.Mhz),
		})
	}
	table.Render()

	color.New(color.FgGreen).Println("\nMemory Information:")
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Total", "Available", "Used", "Used Percent"})
	table.Append([]string{
		fmt.Sprintf("%d", s.VmStat.Total),
		fmt.Sprintf("%d", s.VmStat.Available),
		fmt.Sprintf("%d", s.VmStat.Used),
		fmt.Sprintf("%.2f%%", s.VmStat.UsedPercent),
	})
	table.Render()

	color.New(color.FgGreen).Println("\nDisk Information:")
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Path", "Total", "Free", "Used", "Used Percent"})
	table.Append([]string{
		s.DiscStat.Path,
		fmt.Sprintf("%d", s.DiscStat.Total),
		fmt.Sprintf("%d", s.DiscStat.Free),
		fmt.Sprintf("%d", s.DiscStat.Used),
		fmt.Sprintf("%.2f%%", s.DiscStat.UsedPercent),
	})
	table.Render()

	color.New(color.FgGreen).Println("\nNetwork Information:")
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Bytes Sent", "Bytes Received", "Packets Sent", "Packets Received"})
	for _, net := range s.NetStat {
		table.Append([]string{
			net.Name,
			fmt.Sprintf("%d", net.BytesSent),
			fmt.Sprintf("%d", net.BytesRecv),
			fmt.Sprintf("%d", net.PacketsSent),
			fmt.Sprintf("%d", net.PacketsRecv),
		})
	}
	table.Render()

	color.New(color.FgGreen).Println("\nGo Version:")
	color.New(color.FgYellow).Printf("%s\n", runtime.Version())

}

func SaveSystemInfoToFile(s *SysInfo) {
	file, err := os.Create("system-info.json")
	if err != nil {
		log.Fatalf("Error creating JSON file: %v", err)
		return
	}

	defer file.Close()

	jsonData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
		return
	}

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing JSON data to file: %v", err)
		return
	}

	color.New(color.FgGreen).Printf("\nJSON data saved to system-info.json\n")
}
