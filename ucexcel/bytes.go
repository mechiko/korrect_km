package ucexcel

import (
	"bytes"
)

func (ue *ucexcel) BytesBuffer() (*bytes.Buffer, error) {
	return ue.file.WriteToBuffer()
}
