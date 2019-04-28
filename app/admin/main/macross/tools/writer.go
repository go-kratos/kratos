package tools

import (
	"fmt"
	"go-common/library/log"
	"os"
	"path/filepath"
	"time"
)

type zipSections struct {
	beforeSigningBlock []byte
	signingBlock       []byte
	signingBlockOffset int64
	centraDir          []byte
	centralDirOffset   int64
	eocd               []byte
	eocdOffset         int64
}
type transform func(*zipSections) (*zipSections, error)

func (z *zipSections) writeTo(output string, transform transform) (err error) {
	f, err := os.Create(output)
	if err != nil {
		return
	}

	defer f.Close()

	newZip, err := transform(z)
	if err != nil {
		return
	}

	for _, s := range [][]byte{
		newZip.beforeSigningBlock,
		newZip.signingBlock,
		newZip.centraDir,
		newZip.eocd} {
		_, err := f.Write(s)
		if err != nil {
			return err
		}
	}
	return
}

var _debug bool

// GenerateChannelApk generate a specified channel apk
func GenerateChannelApk(out string, channel string, extras map[string]string, input string, force bool, debug bool) (output string, err error) {
	_debug = debug
	if len(input) == 0 {
		err = fmt.Errorf("Error: no input file specified")
		log.Error("no input file specified, error(%v)", err)
		return
	}

	if _, err = os.Stat(input); os.IsNotExist(err) {
		return
	}

	if len(out) == 0 {
		out = filepath.Dir(input)
	} else {
		var fi os.FileInfo

		err = os.MkdirAll(out, 0755)
		if err != nil {
			log.Error("os.MkdirAll(%s) error(%v)", out, err)
			return
		}

		fi, err = os.Stat(out)
		println("error %v", err)
		if os.IsNotExist(err) || !fi.IsDir() {
			err = fmt.Errorf("Error: output %s is neither exist nor a dir", out)
			log.Error("output %s is neither exist nor a dir, error(%v)", out, err)
			return
		}
	}
	if channel == "" {
		err = fmt.Errorf("Error: no channel specified")
		log.Error("no channel specified, error(%v)", err)
		return
	}

	//TODO: add new option for generating new channel from channelled apk
	if c, _ := readChannelInfo(input); len(c.Channel) != 0 {
		err = fmt.Errorf("Error: file %s is registered a channel block %s", filepath.Base(input), c.String())
		log.Error("file %s is registered a channel block %s, error(%v)", filepath.Base(input), c.String(), err)
		return
	}
	var start time.Time
	if _debug {
		start = time.Now()
	}

	z, err := newZipSections(input)
	if err != nil {
		err = fmt.Errorf("Error occurred on parsing apk %s, %s", input, err)
		log.Error("Error occurred on parsing apk %s, error(%v)", input, err)
		return
	}
	name, ext := fileNameAndExt(input)
	output = filepath.Join(out, name+"-"+channel+ext)
	c := ChannelInfo{Channel: channel, Extras: extras}
	err = gen(c, z, output, force)
	if err != nil {
		err = fmt.Errorf("Error occurred on generating channel %s, %s", channel, err)
		log.Error("Error occurred on generating channel %s, error(%v)", input, err)
		return
	}
	if _debug {
		println("Consume", time.Since(start).String())
	} else {
		println("Done!")
	}
	return
}

func newZipSections(input string) (z zipSections, err error) {
	in, err := os.Open(input)
	if err != nil {
		return
	}
	defer in.Close()

	// read eocd
	eocd, eocdOffset, err := findEndOfCentralDirectoryRecord(in)
	if err != nil {
		return
	}
	centralDirOffset := getEocdCentralDirectoryOffset(eocd)
	centralDirSize := getEocdCentralDirectorySize(eocd)
	z.eocd = eocd
	z.eocdOffset = eocdOffset
	z.centralDirOffset = int64(centralDirOffset)

	// read signing block
	signingBlock, signingBlockOffset, err := findApkSigningBlock(in, centralDirOffset)
	if err != nil {
		return
	}
	z.signingBlock = signingBlock
	z.signingBlockOffset = signingBlockOffset
	// read bytes before signing block
	//TODO: waste too large memory
	if signingBlockOffset >= 64*1024*1024 {
		fmt.Print("Warning: maybe waste large memory on processing this apk! ")
		fmt.Println("Before APK Signing Block bytes size is", signingBlockOffset/1024/1024, "MB")
	}
	beforeSigningBlock := make([]byte, signingBlockOffset)
	n, err := in.ReadAt(beforeSigningBlock, 0)
	if err != nil {
		return
	}
	if int64(n) != signingBlockOffset {
		return z, fmt.Errorf("Read bytes count mismatched! Expect %d, but %d", signingBlockOffset, n)
	}
	z.beforeSigningBlock = beforeSigningBlock

	centralDir := make([]byte, centralDirSize)
	n, err = in.ReadAt(centralDir, int64(centralDirOffset))
	if uint32(n) != centralDirSize {
		return z, fmt.Errorf("Read bytes count mismatched! Expect %d, but %d", centralDirSize, n)
	}
	z.centraDir = centralDir
	if _debug {
		fmt.Printf("signingBlockOffset=%d, signingBlockLenth=%d\n"+
			"centralDirOffset=%d, centralDirSize=%d\n"+
			"eocdOffset=%d, eocdLenthe=%d\n",
			signingBlockOffset,
			len(signingBlock),
			centralDirOffset,
			centralDirSize,
			eocdOffset,
			len(eocd))
	}
	return
}

func gen(info ChannelInfo, sections zipSections, output string, force bool) (err error) {

	fi, err := os.Stat(output)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if fi != nil {
		if !force {
			return fmt.Errorf("file already exists %s", output)
		}
		println("Force generating channel", info.Channel)
	}

	var s time.Time
	if _debug {
		s = time.Now()
	}
	err = sections.writeTo(output, newTransform(info))
	if _debug {
		fmt.Printf("    write %s consume %s", output, time.Since(s).String())
		fmt.Println()
	}
	return
}

func newTransform(info ChannelInfo) transform {
	return func(zip *zipSections) (*zipSections, error) {

		newBlock, diffSize, err := makeSigningBlockWithChannelInfo(info, zip.signingBlock)
		if err != nil {
			return nil, err
		}
		newzip := new(zipSections)
		newzip.beforeSigningBlock = zip.beforeSigningBlock
		newzip.signingBlock = newBlock
		newzip.signingBlockOffset = zip.signingBlockOffset
		newzip.centraDir = zip.centraDir
		newzip.centralDirOffset = zip.centralDirOffset
		newzip.eocdOffset = zip.eocdOffset
		newzip.eocd = makeEocd(zip.eocd, uint32(int64(diffSize)+zip.centralDirOffset))
		return newzip, nil
	}
}
