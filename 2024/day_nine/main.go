package main

import (
	"adventofcode"
	"adventofcode/toolbox/fs"
	"fmt"
)

func main() {
	adventofcode.Time(func() {
		data := fs.LoadFile("2024/day_nine/input.txt")
		disk := LoadDisk(data)

		fmt.Printf("Before Checksum: %d\n", disk.Checksum())
		disk.DefragmentV2()
		fmt.Printf("After Checksum : %d\n", disk.Checksum())
	})
}

type FileBlock struct {
	ID   uint // block IDs start at 1, subtract 1 during calc
	Size uint
}

type Disk []FileBlock

func (d Disk) String() string {
	var s string
	for _, block := range d {
		if block.ID == 0 {
			s += "."
		} else {
			v := fmt.Sprintf("%d", block.ID-1)

			if len(v) >= 2 {
				v = fmt.Sprintf("[%d]", block.ID-1)
			}

			s += v
		}
	}
	return s
}

func (d Disk) DefragmentV1() bool {
	copyHead := uint(len(d)) - 1

	for i := uint(0); i < uint(len(d)); i++ {
		if d[i].ID == 0 {
			for j := copyHead; j > i; j-- {
				if d[j].ID != 0 {
					d[i].ID = d[j].ID
					d[j].ID = 0
					copyHead = j
					break
				}
			}
		}
	}

	return true
}

func (d Disk) DefragmentV2() bool {
	copyHead := int(len(d)) - 1

	// start at the end, seek forward to find a location
	for j := copyHead; j > 0; j-- {
		if d[j].ID != 0 {
			copyHead = j

			fileSize := d[j].Size

			var blankSize uint
			for i := 0; i < j; i++ {
				if d[i].ID != 0 {
					blankSize = 0
				} else {
					blankSize++
				}

				if blankSize == fileSize {
					// move current file to blank
					// i will be last blank and j will be last file
					for k := range int(fileSize) {
						d[i-k] = d[j-k] // might need adjusting
						d[j-k] = FileBlock{ID: 0, Size: 0}
					}
					break
				}
			}
			j -= int(fileSize) - 1
		}
	}

	return true
}

func (d Disk) Checksum() int {
	var checksum uint
	for i, block := range d {
		if block.ID == 0 {
			continue
		}
		checksum += uint(i) * (block.ID - 1)
	}
	return int(checksum)
}

func LoadDisk(data []byte) Disk {
	// understand the disk size
	diskSize := DiskSize(data)
	disk := make([]FileBlock, diskSize)

	id := uint(1)
	writeHead := 0
	for i := uint(0); i < uint(len(data)); i += 2 {
		fileBlock := SingleByteToInt(data[i])

		currentID := id
		for range fileBlock {
			disk[writeHead].ID = currentID
			disk[writeHead].Size = fileBlock

			if currentID == id {
				id++
			}
			writeHead++
		}

		// empty block is not required at end
		if i+1 >= uint(len(data)) {
			break
		}
		emptyBlock := SingleByteToInt(data[i+1])
		for range emptyBlock {
			writeHead++
		}
	}

	return disk
}

func DiskSize(data []byte) uint {
	var diskSize uint
	for _, b := range data {
		diskSize += SingleByteToInt(b)
	}
	return diskSize
}

func SingleByteToInt(data byte) uint {
	var diskSize uint
	switch data {
	case '1':
		diskSize = 1
	case '2':
		diskSize = 2
	case '3':
		diskSize = 3
	case '4':
		diskSize = 4
	case '5':
		diskSize = 5
	case '6':
		diskSize = 6
	case '7':
		diskSize = 7
	case '8':
		diskSize = 8
	case '9':
		diskSize = 9
	}
	return diskSize
}
