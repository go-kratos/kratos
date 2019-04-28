package padding

import "errors"

var (
	ErrPaddingSize = errors.New("pkcs5 padding size error")
	PKCS5          = &pkcs5{}
)

// pkcs5Padding is a pkcs5 padding struct.
type pkcs5 struct{}

// Padding is interface used for crypto.
type Padding interface {
	Padding(src []byte, blockSize int) []byte
	Unpadding(src []byte, blockSize int) ([]byte, error)
}

// Padding implements the Padding interface Padding method.
func (p *pkcs5) Padding(src []byte, blockSize int) []byte {
	srcLen := len(src)
	padLen := byte(blockSize - (srcLen % blockSize))
	pd := make([]byte, srcLen+int(padLen))
	copy(pd, src)
	for i := srcLen; i < len(pd); i++ {
		pd[i] = padLen
	}
	return pd
}

// Unpadding implements the Padding interface Unpadding method.
func (p *pkcs5) Unpadding(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])
	if paddingLen >= srcLen || paddingLen > blockSize {
		return nil, ErrPaddingSize
	}
	return src[:srcLen-paddingLen], nil
}
