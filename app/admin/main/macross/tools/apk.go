package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
)

const (
	apkSigBlockMinSize uint32 = 32

	// https://android.googlesource.com/platform/build/+/android-7.1.2_r27/tools/signapk/src/com/android/signapk/ApkSignerV2.java
	// APK_SIGNING_BLOCK_MAGIC = {
	//  0x41, 0x50, 0x4b, 0x20, 0x53, 0x69, 0x67, 0x20,
	//  0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x20, 0x34, 0x32 }

	apkSigBlockMagicHi = 0x3234206b636f6c42 // LITTLE_ENDIAN, High
	apkSigBlockMagicLo = 0x20676953204b5041 // LITTLE_ENDIAN, Low
	apkChannelBlockID  = 0x71777777
	// https://en.wikipedia.org/wiki/Zip_(file_format)
	// https://android.googlesource.com/platform/build/+/android-7.1.2_r27/tools/signapk/src/com/android/signapk/ZipUtils.java
	zipEocdRecSig                      = 0x06054b50
	zipEocdRecMinSize                  = 22
	zipEocdCentralDirSizeFieldOffset   = 12
	zipEocdCentralDirOffsetFieldOffset = 16
	zipEocdCommentLengthFieldOffset    = 20
)

// ChannelInfo for apk
type ChannelInfo struct {
	Channel string
	Extras  map[string]string
	raw     []byte
}

// ChannelInfo to string
func (c *ChannelInfo) String() string {
	b := c.Bytes()
	if b == nil {
		return ""
	}
	return string(b)
}

// Bytes for ChannelInfo to byte array
func (c *ChannelInfo) Bytes() []byte {
	if c.raw != nil {
		return c.raw
	}
	if len(c.Channel) == 0 && c.Extras == nil {
		return nil
	}
	var buf bytes.Buffer
	buf.WriteByte('{')
	if len(c.Channel) != 0 {
		buf.WriteString("\"channel\":")
		buf.WriteByte('"')
		buf.WriteString(c.Channel)
		buf.WriteByte('"')
		buf.WriteByte(',')
	}

	if c.Extras != nil {
		for k, v := range c.Extras {
			buf.WriteByte('"')
			buf.WriteString(k)
			buf.WriteByte('"')
			buf.WriteByte(':')
			buf.WriteByte('"')
			buf.WriteString(v)
			buf.WriteByte('"')
			buf.WriteByte(',')
		}
	}
	if buf.Len() > 2 {
		buf.Truncate(buf.Len() - 1)
	}

	buf.WriteByte('}')

	return buf.Bytes()
}

func readChannelInfo(file string) (c ChannelInfo, err error) {
	block, err := readChannelBlock(file)
	if err != nil {
		return c, err
	}

	if block != nil {
		var bundle map[string]string
		err := json.Unmarshal(block, &bundle)
		if err != nil {
			return c, err
		}
		c.Channel = bundle["channel"]
		delete(bundle, "channel")
		c.Extras = bundle
		c.raw = block
	}
	return c, nil
}

// read block associated to apkChannelBlockID
func readChannelBlock(file string) ([]byte, error) {
	m, err := readIDValues(file, apkChannelBlockID)
	if err != nil {
		return nil, err
	}
	return m[apkChannelBlockID], nil
}

func readIDValues(file string, ids ...uint32) (map[uint32][]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	eocd, offset, err := findEndOfCentralDirectoryRecord(f)
	if err != nil {
		return nil, err
	}
	if offset <= 0 {
		return nil, errors.New("Cannot find EOCD record, maybe a broken zip file")
	}
	centralDirOffset := getEocdCentralDirectoryOffset(eocd)
	block, _, err := findApkSigningBlock(f, centralDirOffset)
	if err != nil {
		return nil, err
	}
	return findIDValuesInApkSigningBlock(block, ids...)
}

// End of central directory record (EOCD)
//
// Offset    Bytes     Description[23]
// 0           4       End of central directory signature = 0x06054b50
// 4           2       Number of this disk
// 6           2       Disk where central directory starts
// 8           2       Number of central directory records on this disk
// 10          2       Total number of central directory records
// 12          4       Size of central directory (bytes)
// 16          4       Offset of start of central directory, relative to start of archive
// 20          2       Comment length (n)
// 22          n       Comment
// For a zip with no archive comment, the
// end-of-central-directory record will be 22 bytes long, so
// we expect to find the EOCD marker 22 bytes from the end.
func findEndOfCentralDirectoryRecord(f *os.File) ([]byte, int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, -1, err
	}
	if fi.Size() < zipEocdRecMinSize {
		// No space for EoCD record in the file.
		return nil, -1, nil
	}
	// Optimization: 99.99% of APKs have a zero-length comment field in the EoCD record and thus
	// the EoCD record offset is known in advance. Try that offset first to avoid unnecessarily
	// reading more data.
	ret, offset, err := findEOCDRecord(f, 0)
	if err != nil {
		return nil, -1, err
	}
	if ret != nil && offset != -1 {
		return ret, offset, nil
	}
	// EoCD does not start where we expected it to. Perhaps it contains a non-empty comment
	// field. Expand the search. The maximum size of the comment field in EoCD is 65535 because
	// the comment length field is an unsigned 16-bit number.
	return findEOCDRecord(f, math.MaxUint16)
}

func findEOCDRecord(f *os.File, maxCommentSize uint16) ([]byte, int64, error) {
	if maxCommentSize > math.MaxUint16 {
		return nil, -1, os.ErrInvalid
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, -1, err
	}
	fileSize := fi.Size()
	if fileSize < zipEocdRecMinSize {
		// No space for EoCD record in the file.
		return nil, -1, nil
	}
	// Lower maxCommentSize if the file is too small.
	if s := uint16(fileSize - zipEocdRecMinSize); maxCommentSize > s {
		maxCommentSize = s
	}
	maxEocdSize := zipEocdRecMinSize + maxCommentSize
	bufOffsetInFile := fileSize - int64(maxEocdSize)
	buf := make([]byte, maxEocdSize)
	n, e := f.ReadAt(buf, bufOffsetInFile)
	if e != nil {
		return nil, -1, err
	}
	eocdOffsetInFile :=
		func() int64 {
			eocdWithEmptyCommentStartPosition := n - zipEocdRecMinSize
			for expectedCommentLength := uint16(0); expectedCommentLength < maxCommentSize; expectedCommentLength++ {
				eocdStartPos := eocdWithEmptyCommentStartPosition - int(expectedCommentLength)
				if getUint32(buf, eocdStartPos) == zipEocdRecSig {
					n := eocdStartPos + zipEocdCommentLengthFieldOffset
					actualCommentLength := getUint16(buf, n)
					if actualCommentLength == expectedCommentLength {
						return int64(eocdStartPos)
					}
				}
			}
			return -1
		}()
	if eocdOffsetInFile == -1 {
		// No EoCD record found in the buffer
		return nil, -1, nil
	}
	// EoCD found
	return buf[eocdOffsetInFile:], bufOffsetInFile + eocdOffsetInFile, nil

}

func getEocdCentralDirectoryOffset(buf []byte) uint32 {
	return getUint32(buf, zipEocdCentralDirOffsetFieldOffset)
}
func getEocdCentralDirectorySize(buf []byte) uint32 {
	return getUint32(buf, zipEocdCentralDirSizeFieldOffset)
}

func setEocdCentralDirectoryOffset(eocd []byte, offset uint32) {
	putUint32(offset, eocd, zipEocdCentralDirOffsetFieldOffset)
}

func isExpected(ids []uint32, test uint32) bool {
	for _, id := range ids {
		if id == test {
			return true
		}
	}
	return false
}

func findIDValuesInApkSigningBlock(block []byte, ids ...uint32) (map[uint32][]byte, error) {
	ret := make(map[uint32][]byte)
	position := 8
	limit := len(block) - 24
	entryCount := 0
	for limit > position { // has remaining bytes
		entryCount++
		if limit-position < 8 { // but not enough
			return nil, fmt.Errorf("APK Signing Block broken on entry #%d", entryCount)
		}
		length := int(getUint64(block, position))
		position += 8

		if length < 4 || length > limit-position {
			return nil, fmt.Errorf("APK Signing Block broken on entry #%d,"+
				" size out of range: length=%d, remaining=%d", entryCount, length, limit-position)
		}
		nextEntryPosition := position + length
		id := getUint32(block, position)
		position += 4
		if len(ids) == 0 || isExpected(ids, id) {
			ret[id] = block[position : position+length-4]
		}
		position = nextEntryPosition

	}
	return ret, nil
}

// Find the APK Signing Block. The block immediately precedes the Central Directory.
//
// FORMAT:
//	 uint64:  size (excluding this field)
//	 repeated ID-value pairs:
//	     uint64:           size (excluding this field)
//	     uint32:           ID
//	     (size - 4) bytes: value
//	 uint64:  size (same as the one above)
//	 uint128: magic
func findApkSigningBlock(f *os.File, centralDirOffset uint32) (block []byte, offset int64, err error) {

	if centralDirOffset < apkSigBlockMinSize {
		return block, offset, fmt.Errorf("APK too small for APK Signing Block."+
			" ZIP Central Directory offset: %d", centralDirOffset)
	}
	// Read the footer of APK signing block
	// 24 = sizeof(uint128) + sizeof(uint64)
	footer := make([]byte, 24)
	_, err = f.ReadAt(footer, int64(centralDirOffset-24))
	if err != nil {
		return
	}
	// Read the magic and block size
	var blockSizeInFooter = getUint64(footer, 0)
	if blockSizeInFooter < 24 || blockSizeInFooter > uint64(math.MaxInt32-8 /* ID-value size field*/) {
		return block, offset, fmt.Errorf("APK Signing Block size out of range: %d", blockSizeInFooter)
	}
	if getUint64(footer, 8) != apkSigBlockMagicLo ||
		getUint64(footer, 16) != apkSigBlockMagicHi {
		return block, offset, errors.New("No APK Signing Block before ZIP Central Directory")
	}

	totalSize := blockSizeInFooter + 8 /* APK signing block size field*/

	offset = int64(uint64(centralDirOffset) - totalSize)

	if offset <= 0 {
		return block, offset, fmt.Errorf("invalid offset for APK Signing Block %d", offset)
	}
	block = make([]byte, totalSize)
	_, err = f.ReadAt(block, offset)
	if err != nil {
		return
	}
	blockSizeInHeader := getUint64(block, 0)
	if blockSizeInHeader != blockSizeInFooter {
		return nil, offset, fmt.Errorf("APK Signing Block sizes in header "+
			"and footer are mismatched! Except %d but %d", blockSizeInFooter, blockSizeInHeader)
	}

	return block, offset, nil
}

// FORMAT:
// uint64:  size (excluding this field)
// repeated ID-value pairs:
//     uint64:           size (excluding this field)
//     uint32:           ID
//     (size - 4) bytes: value
// uint64:  size (same as the one above)
// uint128: magic
func makeSigningBlockWithChannelInfo(info ChannelInfo, signingBlock []byte) ([]byte, int, error) {

	signingBlockSize := getUint64(signingBlock, 0)
	signingBlockLen := len(signingBlock)
	if n := uint64(signingBlockLen - 8); signingBlockSize != n {
		return nil, 0, fmt.Errorf("APK Signing Block is illegal! Expect size %d but %d", signingBlockSize, n)
	}

	channelValue := info.Bytes()
	channelValueSize := uint64(4 + len(channelValue))
	resultSize := 8 + signingBlockSize + 8 + channelValueSize

	newBlock := make([]byte, resultSize)
	position := 0
	putUint64(resultSize-8, newBlock, position)
	position += 8
	// copy raw id-value pairs
	n, _ := copyBytes(signingBlock, 8, newBlock, position, int(signingBlockSize)-16-8)
	position += n
	putUint64(channelValueSize, newBlock, position)
	position += 8
	putUint32(apkChannelBlockID, newBlock, position)
	position += 4
	n, _ = copyBytes(channelValue, 0, newBlock, position, len(channelValue))
	position += n

	putUint64(resultSize-8, newBlock, position)
	position += 8
	copyBytes(signingBlock, signingBlockLen-16, newBlock, int(resultSize-16), 16)
	position += 16

	if position != int(resultSize) {
		return nil, -1, fmt.Errorf("count mismatched ! %d vs %d", position, resultSize)
	}
	return newBlock, int(resultSize) - signingBlockLen, nil
}

func makeEocd(origin []byte, newCentralDirOffset uint32) []byte {
	eocd := make([]byte, len(origin))
	copy(eocd, origin)
	setEocdCentralDirectoryOffset(eocd, newCentralDirOffset)
	return eocd
}
