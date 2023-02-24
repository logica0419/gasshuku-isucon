package validator

import (
	"errors"
	"fmt"
	"image/png"
	"net/http"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

const (
	qrCodeSize         = 49
	errCorrectionLevel = "M"
)

// QRコードの内容が一致するか検証する
func WithQRCodeEqual(content string, decryptFunc func(string) (string, error)) ValidateOpt {
	return func(res *http.Response) error {
		img, err := png.Decode(res.Body)
		if err != nil {
			return failure.NewError(model.ErrUndecodableBody, err)
		}

		rect := img.Bounds()
		if rect.Dx() != qrCodeSize || rect.Dy() != qrCodeSize {
			return failure.NewError(model.ErrInvalidBody, errors.New("size is not 45x45"))
		}

		bmp, err := gozxing.NewBinaryBitmapFromImage(img)
		if err != nil {
			return failure.NewError(model.ErrUndecodableBody, err)
		}

		qrReader := qrcode.NewQRCodeReader()
		result, err := qrReader.Decode(bmp, nil)
		if err != nil {
			return failure.NewError(model.ErrUndecodableBody, err)
		}

		if result.GetResultMetadata()[gozxing.ResultMetadataType_ERROR_CORRECTION_LEVEL] != errCorrectionLevel {
			return failure.NewError(model.ErrInvalidBody, fmt.Errorf("invalid error correction level: expected: M, actual: %s", result.GetResultMetadata()[gozxing.ResultMetadataType_ERROR_CORRECTION_LEVEL].(string)))
		}

		decryptedContent, err := decryptFunc(result.String())
		if err != nil {
			return failure.NewError(model.ErrInvalidBody, err)
		}

		if decryptedContent != content {
			return failure.NewError(model.ErrInvalidBody, fmt.Errorf("invalid content: expected %s, actual %s", content, result.GetText()))
		}

		return nil
	}
}
