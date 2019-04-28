// +build linux darwin

package resolvconf

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const (
	resolvConfPath = "/etc/resolv.conf"
)

// ParseResolvConf parse /etc/resolv.conf file and return nameservers
func ParseResolvConf() ([]string, error) {
	fp, err := os.Open(resolvConfPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return parse(fp)
}

func parse(fp io.Reader) ([]string, error) {
	var result []string

	bufRd := bufio.NewReader(fp)
	for {
		line, err := bufRd.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			if line == "" {
				break
			}
		}
		line = strings.TrimSpace(line)

		// ignore comment, comment startwith #
		if strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if fields[0] == "nameserver" {
			result = append(result, fields[1:]...)
		}
	}
	return result, nil
}
