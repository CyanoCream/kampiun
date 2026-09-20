package pgdb

import (
	"encoding/json"
	"fmt"

	"kampiun/domain"
)

// sportFormatBytes menyerakkan format jadi JSON BLOB.
func sportFormatBytes(f domain.SportFormat) ([]byte, error) {
	return json.Marshal(f)
}

func sportFormatFromBytes(b []byte) domain.SportFormat {
	var f domain.SportFormat
	if len(b) > 0 {
		_ = json.Unmarshal(b, &f)
	}
	return f
}

// fmtAr = fmt.Sprintf yang menerima 1..n argumen angka posisi.
func fmtAr(format string, n ...int) string {
	args := make([]any, len(n))
	for i, v := range n {
		args[i] = v
	}
	return fmt.Sprintf(format, args...)
}
