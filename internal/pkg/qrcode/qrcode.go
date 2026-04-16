package qrcode

import (
	"hotel-loyalty/internal/infrastructure/logger"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(content string) ([]byte, error) {
	logger.InfoLogger.Println(
		"[QRCODE] generating",
		"content=", content,
	)

	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		logger.ErrorLogger.Println(
			"[QRCODE] generate failed",
			"error=", err,
			"content=", content,
		)
		return nil, err
	}

	logger.InfoLogger.Println(
		"[QRCODE] generated successfully",
	)

	return png, nil
}
