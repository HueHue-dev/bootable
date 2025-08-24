package disk

import "os"

type linuxDisk struct {
	*os.File
}

func openDisk(path string) (Disk, error) {
	f, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &linuxDisk{File: f}, nil
}
