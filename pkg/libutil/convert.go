package libutil

import "strconv"

// convert human readable size to bytes
func ConvertSizeToBytes(size string) (int64, error) {
	var sizeInBytes int64
	var err error

	switch size[len(size)-1] {
	case 'k', 'K':
		sizeInBytes, err = strconv.ParseInt(size[:len(size)-1], 10, 64)
		if err != nil {
			return 0, err
		}
		sizeInBytes *= 1024
	case 'm', 'M':
		sizeInBytes, err = strconv.ParseInt(size[:len(size)-1], 10, 64)
		if err != nil {
			return 0, err
		}
		sizeInBytes *= 1024 * 1024
	case 'g', 'G':
		sizeInBytes, err = strconv.ParseInt(size[:len(size)-1], 10, 64)
		if err != nil {
			return 0, err
		}
		sizeInBytes *= 1024 * 1024 * 1024
	default:
		sizeInBytes, err = strconv.ParseInt(size, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	return sizeInBytes, nil
}

// convert bytes to human readable size
func ConvertBytesToSize(size int64) string {
	var sizeInString string

	switch {
	case size >= 1024*1024*1024:
		sizeInString = strconv.FormatInt(size/(1024*1024*1024), 10) + "G"
	case size >= 1024*1024:
		sizeInString = strconv.FormatInt(size/(1024*1024), 10) + "M"
	case size >= 1024:
		sizeInString = strconv.FormatInt(size/1024, 10) + "K"
	default:
		sizeInString = strconv.FormatInt(size, 10)
	}

	return sizeInString
}
