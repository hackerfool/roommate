package api

import (
	"crypto/aes"
	"crypto/cipher"
	"mlog"
)

func aes128cbc(key []byte, data []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	data = PKCS7UnPadding(data)
	orgdata := make([]byte, len(data))
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(orgdata, data)

	return PKCS7UnPadding(orgdata), nil
}

func pkcs7(bsize int, data []byte) []byte {
	var padding int
	r := len(data) % bsize
	if r == 0 {
		padding = bsize
	} else {
		padding = r
	}
	for i := 0; i < padding; i++ {
		data = append(data, byte(padding))
	}
	mlog.Debug(data)

	return data
}

// PKCS7UnPadding return unpadding []Byte plantText
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return plantText[:(length - unPadding)]
}
