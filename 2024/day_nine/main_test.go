package main

import (
	"fmt"
	"testing"
)

func TestDiskSize(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"1254", 12},
		{"2333133121414131402", 42},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			diskSize := DiskSize([]byte(test.in))
			if int(diskSize) != test.want {
				t.Errorf("disk size not as expected, got: %d, want: %d", diskSize, test.want)
			}
		})
	}
}

func TestLoadDiskV1(t *testing.T) {
	tests := []struct {
		in       string
		want     int
		fragged  string
		checksum int
	}{
		{"1254", 12, "011111......", 15},
		{"12345", 15, "022111222......", 60},
		{"2333133121414131402", 42, "0099811188827773336446555566..............", 1928},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			disk := LoadDisk([]byte(test.in))

			fmt.Printf("%s\n", disk)
			disk.DefragmentV1()
			fmt.Printf("%s\n", disk)
			if fmt.Sprintf("%s", disk) != test.fragged {
				t.Errorf("disk not as expected, got: %s, want: %s", fmt.Sprintf("%s", disk), test.fragged)
			}

			if disk.Checksum() != test.checksum {
				t.Errorf("checksum not as expected, got: %d, want: %d", disk.Checksum(), test.checksum)
			}
		})
	}
}

func TestLoadDiskV2(t *testing.T) {
	tests := []struct {
		in       string
		want     int
		fragged  string
		checksum int
	}{
		{"1254", 12, "0..11111....", 25},
		{"2333133121414131402", 42, "00992111777.44.333....5555.6666.....8888..", 2858},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			disk := LoadDisk([]byte(test.in))

			fmt.Printf("%s\n", disk)
			disk.DefragmentV2()
			fmt.Printf("%s\n", disk)
			if fmt.Sprintf("%s", disk) != test.fragged {
				t.Errorf("disk not as expected, got: %s, want: %s", fmt.Sprintf("%s", disk), test.fragged)
			}

			if disk.Checksum() != test.checksum {
				t.Errorf("checksum not as expected, got: %d, want: %d", disk.Checksum(), test.checksum)
			}
		})
	}
}
