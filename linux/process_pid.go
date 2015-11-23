package linux

import (
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

func ReadMaxPID(path string) (uint64, error) {

	b, err := ioutil.ReadFile(path)

	if err != nil {
		return 0, err
	}

	s := strings.TrimSpace(string(b))

	i, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		return 0, err
	}

	return i, nil

}

type UInt64Slice []uint64

func (p UInt64Slice) Len() int           { return len(p) }
func (p UInt64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p UInt64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func ListPID(path string) ([]uint64, error) {
	pids := make([]uint64, 0, 5)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return pids, err
	}
	for _, file := range files {
		if file.IsDir() {
			if name, err := strconv.Atoi(file.Name()); err == nil {
				pids = append(pids, uint64(name))
			}
		}
	}
	sort.Sort(UInt64Slice(pids))
	return pids, nil
}
