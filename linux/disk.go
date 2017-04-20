package linux

import (
	"syscall"
)

type Disk struct {
	Size       uint64 `json:"size"`
	SizeUsed   uint64 `json:"size_used"`
	SizeFree   uint64 `json:"size_free"`
	SizeAvail  uint64 `json:"size_avail"`
	Inodes     uint64 `json:"inodes"`
	InodesUsed uint64 `json:"inodes_used"`
	InodesFree uint64 `json:"inodes_free"`
}

func ReadDisk(path string) (*Disk, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return nil, err
	}
	disk := Disk{}
	disk.Size = fs.Blocks * uint64(fs.Bsize)
	disk.SizeFree = fs.Bfree * uint64(fs.Bsize)
	disk.SizeUsed = disk.Size - disk.SizeFree
	disk.SizeAvail = fs.Bavail * uint64(fs.Bsize)
	disk.Inodes = fs.Files
	disk.InodesFree = fs.Ffree
	disk.InodesUsed = fs.Files - fs.Ffree
	return &disk, nil
}
