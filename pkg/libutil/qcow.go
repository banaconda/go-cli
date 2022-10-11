package libutil

import (
	"os"
	"os/exec"
	"strconv"
)

// check is qcow2 file
func IsQcow2(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	file.Read(buffer)
	return string(buffer[0:4]) == "QFI\xfb"
}

// resize qcow2 image
func ResizeQcow2Image(imagePath string, size int64) (string, error) {
	// resize image
	cmd := exec.Command("qemu-img", "resize", imagePath, strconv.FormatInt(size, 10))
	out, err := cmd.Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}
