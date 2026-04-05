package base

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var protobufRawDescBlockRE = regexp.MustCompile(`(?ms)^const file_.*?_rawDesc = "" \+\r?\n(?:\t"(?:[^"\\]|\\.)*" \+\r?\n)*\t"(?:[^"\\]|\\.)*"\r?\n`)

func kratosHome() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	home := filepath.Join(dir, ".kratos")
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0o700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func kratosHomeWithDir(dir string) string {
	home := filepath.Join(kratosHome(), dir)
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0o700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func copyFile(src, dst string, replaces []string) error {
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	buf, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	buf = replaceTemplateContent(buf, replaces)
	return os.WriteFile(dst, buf, srcinfo.Mode())
}

func replaceTemplateContent(buf []byte, replaces []string) []byte {
	matches := protobufRawDescBlockRE.FindAllIndex(buf, -1)
	if len(matches) == 0 {
		return applyReplacements(buf, replaces)
	}

	var out bytes.Buffer
	last := 0
	for _, match := range matches {
		out.Write(applyReplacements(buf[last:match[0]], replaces))
		out.Write(buf[match[0]:match[1]])
		last = match[1]
	}
	out.Write(applyReplacements(buf[last:], replaces))
	return out.Bytes()
}

func applyReplacements(buf []byte, replaces []string) []byte {
	for i := 0; i+1 < len(replaces); i += 2 {
		buf = bytes.ReplaceAll(buf, []byte(replaces[i]), []byte(replaces[i+1]))
	}
	return buf
}

func copyDir(src, dst string, replaces, ignores []string) error {
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcinfo.Mode())
	if err != nil {
		return err
	}

	fds, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, fd := range fds {
		if hasSets(fd.Name(), ignores) {
			continue
		}
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())
		var e error
		if fd.IsDir() {
			e = copyDir(srcfp, dstfp, replaces, ignores)
		} else {
			e = copyFile(srcfp, dstfp, replaces)
		}
		if e != nil {
			return e
		}
	}
	return nil
}

func hasSets(name string, sets []string) bool {
	for _, ig := range sets {
		if ig == name {
			return true
		}
	}
	return false
}

func Tree(path string, dir string) {
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			fmt.Printf("%s %s (%v bytes)\n", color.GreenString("CREATED"), strings.ReplaceAll(path, dir+"/", ""), info.Size())
		}
		return nil
	})
}
