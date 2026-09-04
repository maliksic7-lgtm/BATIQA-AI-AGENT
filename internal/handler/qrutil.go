package handler

import (
	"bytes"

	qrcode "github.com/skip2/go-qrcode"
)

// qrcodeEncode renders text as PNG QR of the given pixel size.
func qrcodeEncode(text string, size int) ([]byte, error) {
	buf := new(bytes.Buffer)
	q, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	png, err := q.PNG(size)
	if err != nil {
		return nil, err
	}
	buf.Write(png)
	return buf.Bytes(), nil
}
