package qrcode

import (
	"bytes"
	"image/png"
	
	"github.com/skip2/go-qrcode"
)

func GenerateQR(data string, size int) ([]byte, error) {
	qr, err := qrcode.Encode(data, qrcode.Medium, size)
	if err != nil {
		return nil, err
	}
	
	// qr уже []byte в формате PNG
	return qr, nil
}
