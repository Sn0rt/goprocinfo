package linux

import (
	"syscall"
)

type Disk struct {
	Inodes     uint64 `json:"inodes"`
	InodesUsed uint64 `json:"inodes_used"`
	InodesFree uint64 `json:"inodes_free"`
	Blocks     uint64 `json:"blocks"`
	BlocksFree uint64 `json:"blocks_free"`
	BlocksUsed uint64 `json:"blocks_used"`
	Size       uint64 `json:"size"`
	SizeUsed   uint64 `json:"size_used"`
	SizeFree   uint64 `json:"size_free"`
	SizeAvail  uint64 `json:"size_avail"`
}

// Bavail is the amount of blocks available to unprivileged users.
// Bfree is simply the total number of free blocks

func ReadDisk(path string) (*Disk, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return nil, err
	}
	disk := Disk{}
	disk.Inodes = fs.Files
	disk.InodesFree = fs.Ffree
	disk.InodesUsed = fs.Files - fs.Ffree
	disk.Blocks = fs.Blocks
	disk.BlocksUsed = fs.Blocks - fs.Bfree
	disk.BlocksFree = fs.Bfree
	disk.Size = fs.Blocks * uint64(fs.Bsize)
	disk.SizeFree = fs.Bfree * uint64(fs.Bsize)
	disk.SizeUsed = disk.Size - disk.SizeFree
	disk.SizeAvail = fs.Bavail * uint64(fs.Bsize)
	return &disk, nil
}
