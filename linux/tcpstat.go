package linux

import (
	"fmt"
	"io/ioutil"
	"os"
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
	State string
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
		state := TcpState[line_array[3]]
		tcpstat := TcpStat{state}

		tcpstats = append(tcpstats, tcpstat)
	}

	return tcpstats, nil
}
