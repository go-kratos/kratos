package cpu

import (
	"bufio"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func readFile(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "os/stat: read file(%s) failed!", path)
	}
	return strings.TrimSpace(string(contents)), nil
}

func parseUint(s string) (uint64, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, 10, 64)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil &&
			intErr.(*strconv.NumError).Err == strconv.ErrRange &&
			intValue < 0 {
			return 0, nil
		}
		return 0, errors.Wrapf(err, "os/stat: parseUint(%s) failed!", s)
	}
	return v, nil
}

// ParseUintList parses and validates the specified string as the value
// found in some cgroup file (e.g. cpuset.cpus, cpuset.mems), which could be
// one of the formats below. Note that duplicates are actually allowed in the
// input string. It returns a map[int]bool with available elements from val
// set to true.
// Supported formats:
// 7
// 1-6
// 0,3-4,7,8-10
// 0-0,0,1-7
// 03,1-3 <- this is gonna get parsed as [1,2,3]
// 3,2,1
// 0-2,3,1
func ParseUintList(val string) (map[int]bool, error) {
	if val == "" {
		return map[int]bool{}, nil
	}

	availableInts := make(map[int]bool)
	split := strings.Split(val, ",")
	errInvalidFormat := errors.Errorf("os/stat: invalid format: %s", val)
	for _, r := range split {
		if !strings.Contains(r, "-") {
			v, err := strconv.Atoi(r)
			if err != nil {
				return nil, errInvalidFormat
			}
			availableInts[v] = true
		} else {
			split := strings.SplitN(r, "-", 2)
			min, err := strconv.Atoi(split[0])
			if err != nil {
				return nil, errInvalidFormat
			}
			max, err := strconv.Atoi(split[1])
			if err != nil {
				return nil, errInvalidFormat
			}
			if max < min {
				return nil, errInvalidFormat
			}
			for i := min; i <= max; i++ {
				availableInts[i] = true
			}
		}
	}
	return availableInts, nil
}

// ReadLines reads contents from a file and splits them by new lines.
// A convenience wrapper to ReadLinesOffsetN(filename, 0, -1).
func readLines(filename string) ([]string, error) {
	return readLinesOffsetN(filename, 0, -1)
}

// ReadLinesOffsetN reads contents from file and splits them by new line.
// The offset tells at which line number to start.
// The count determines the number of lines to read (starting from offset):
//   n >= 0: at most n lines
//   n < 0: whole file
func readLinesOffsetN(filename string, offset uint, n int) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []string{""}, err
	}
	defer f.Close()

	var ret []string

	r := bufio.NewReader(f)
	for i := 0; i < n+int(offset) || n < 0; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if i < int(offset) {
			continue
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}

	return ret, nil
}
