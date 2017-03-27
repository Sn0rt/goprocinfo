package linux

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type SlabStat struct {
	SlabInfos       []SlabInfo `json:"slabs"`
	TotalSize       uint64     `json:"total_size"`
	TotalActiveSize uint64     `json:"total_active_size"`
	NumObjs         uint64     `json:"num_objs"`
	NumActiveObjs   uint64     `json:"num_active_objs"`
	NumPages        uint64     `json:"num_pages"`
	NumSlabs        uint64     `json:"num_slabs"`
	NumActiveSlabs  uint64     `json:"num_active_slabs"`
	NumCaches       uint64     `json:"num_caches"`
	NumActiveCaches uint64     `json:"num_active_caches"`
}

type SlabInfo struct {
	Name         string `json:"name"`
	ActiveObjs   uint64 `json:"active_objs"`
	NumObjs      uint64 `json:"num_objs"`
	ObjSize      uint64 `json:"obj_size"`
	ObjPerSlab   uint64 `json:"obj_per_slab"`
	PagesPerSlab uint64 `json:"pages_per_slab"`
	ActiveSlabs  uint64 `json:"active_slabs"`
	NumSlabs     uint64 `json:"num_slabs"`
	Use          uint64 `json:"use"`
	CacheSize    uint64 `json:"cache_size"`
}

func ReadSlabStat(path string) (*SlabStat, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(b)
	lines := strings.Split(content, "\n")
	pageSize := uint64(os.Getpagesize())

	var slabStat = SlabStat{}
	var slabInfo = SlabInfo{}

	for _, line := range lines[2:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		slabInfo.Name = fields[0]
		slabInfo.ActiveObjs, _ = strconv.ParseUint(fields[1], 10, 64)
		slabInfo.NumObjs, _ = strconv.ParseUint(fields[2], 10, 64)
		slabInfo.ObjSize, _ = strconv.ParseUint(fields[3], 10, 64)
		slabInfo.ObjPerSlab, _ = strconv.ParseUint(fields[4], 10, 64)
		slabInfo.PagesPerSlab, _ = strconv.ParseUint(fields[5], 10, 64)
		slabInfo.ActiveSlabs, _ = strconv.ParseUint(fields[13], 10, 64)
		slabInfo.NumSlabs, _ = strconv.ParseUint(fields[14], 10, 64)
		if slabInfo.NumObjs != 0 {
			slabInfo.Use = 100 * slabInfo.ActiveObjs / slabInfo.NumObjs
			slabStat.NumActiveCaches++
		} else {
			slabInfo.Use = 0
		}
		slabInfo.CacheSize = slabInfo.NumSlabs * slabInfo.PagesPerSlab * pageSize

		slabStat.SlabInfos = append(slabStat.SlabInfos, slabInfo)
		slabStat.TotalSize += slabInfo.NumObjs * slabInfo.ObjSize
		slabStat.TotalActiveSize += slabInfo.ActiveObjs * slabInfo.ObjSize
		slabStat.NumObjs += slabInfo.NumObjs
		slabStat.NumActiveObjs += slabInfo.ActiveObjs
		slabStat.NumPages += slabInfo.NumSlabs * slabInfo.PagesPerSlab
		slabStat.NumSlabs += slabInfo.NumSlabs
		slabStat.NumActiveSlabs += slabInfo.ActiveSlabs
		slabStat.NumCaches++
	}
	return &slabStat, nil
}
