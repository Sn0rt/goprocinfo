package linux

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var TcpState = map[string]string{
	"01": "established",
	"02": "syn_sent",
	"03": "syn_recv",
	"04": "fin_wait1",
	"05": "fin_wait2",
	"06": "time_wait",
	"07": "close",
	"08": "close_wait",
	"09": "last_ack",
	"0A": "listen",
	"0B": "closing",
}

type TcpStat struct {
	State      string
	LocalIp    string
	LocalPort  int64
	RemoteIp   string
	RemotePort int64
}

func getLines(path string) []string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")

	return lines[1 : len(lines)-1]

}

func hexToDec(h string) int64 {
	d, err := strconv.ParseInt(h, 16, 32)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return d
}

func convertIp(ip string) string {
	var out string
	if len(ip) > 8 {
		i := []string{ip[30:32], ip[28:30], ip[26:28], ip[24:26], ip[22:24], ip[20:22], ip[18:20], ip[16:18], ip[14:16], ip[12:14], ip[10:12], ip[8:10], ip[6:8], ip[4:6], ip[2:4], ip[0:2]}
		out = fmt.Sprintf("%v%v:%v%v:%v%v:%v%v:%v%v:%v%v:%v%v:%v%v", i[14], i[15], i[13], i[12], i[10], i[11], i[8], i[9], i[6], i[7], i[4], i[5], i[2], i[3], i[0], i[1])

	} else {
		i := []int64{hexToDec(ip[6:8]), hexToDec(ip[4:6]), hexToDec(ip[2:4]), hexToDec(ip[0:2])}
		out = fmt.Sprintf("%v.%v.%v.%v", i[0], i[1], i[2], i[3])
	}
	return out
}

func removeEmpty(array []string) []string {
	// remove empty data from line
	var new_array []string
	for _, i := range array {
		if i != "" {
			new_array = append(new_array, i)
		}
	}
	return new_array
}

func ReadTcpStats(path string) ([]TcpStat, error) {
	lines := getLines(path)

	var tcpstats []TcpStat

	for _, line := range lines {

		line_array := removeEmpty(strings.Split(strings.TrimSpace(line), " "))
		ip_port := strings.Split(line_array[1], ":")
		ip := convertIp(ip_port[0])
		port := hexToDec(ip_port[1])

		fip_port := strings.Split(line_array[2], ":")
		fip := convertIp(fip_port[0])
		fport := hexToDec(fip_port[1])

		state := TcpState[line_array[3]]
		tcpstat := TcpStat{state, ip, port, fip, fport}

		tcpstats = append(tcpstats, tcpstat)
	}

	return tcpstats, nil
}
