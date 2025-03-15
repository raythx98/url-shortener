package qrcode

import "github.com/skip2/go-qrcode"

type IQrCode interface {
	Encode(content string) ([]byte, error)
}

type QrCode struct{}

func New() *QrCode {
	return &QrCode{}
}

func (q *QrCode) Encode(content string) ([]byte, error) {
	return qrcode.Encode(content, qrcode.Medium, 256)
}
