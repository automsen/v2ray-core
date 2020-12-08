package serial

import (
	"bytes"
	"encoding/json"
	"io"

	"encoding/base64"
	"github.com/forgoer/openssl"
	"v2ray.com/core"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/infra/conf"
	json_reader "v2ray.com/core/infra/conf/json"
)

type offset struct {
	line int
	char int
}

func findOffset(b []byte, o int) *offset {
	if o >= len(b) || o < 0 {
		return nil
	}

	line := 1
	char := 0
	for i, x := range b {
		if i == o {
			break
		}
		if x == '\n' {
			line++
			char = 0
		} else {
			char++
		}
	}

	return &offset{line: line, char: char}
}

func LoadJSONConfig(reader io.Reader, entry bool, entryKey string) (*core.Config, error) {
	jsonConfig := &conf.Config{}

	jsonContent := bytes.NewBuffer(make([]byte, 0, 10240))

	//modify by tanglongsen
	if entry {

		data, err := buf.ReadAllToBytes(reader)
		if err != nil {
			return nil, err
		}

		decodeBytes, _ := base64.StdEncoding.DecodeString(string(data))
		dst, _ := openssl.AesECBDecrypt(decodeBytes, []byte(entryKey), openssl.PKCS7_PADDING)

		reader = bytes.NewReader(dst)

	}

	jsonReader := io.TeeReader(&json_reader.Reader{
		Reader: reader,
	}, jsonContent)
	decoder := json.NewDecoder(jsonReader)

	if err := decoder.Decode(jsonConfig); err != nil {
		var pos *offset
		cause := errors.Cause(err)
		switch tErr := cause.(type) {
		case *json.SyntaxError:
			pos = findOffset(jsonContent.Bytes(), int(tErr.Offset))
		case *json.UnmarshalTypeError:
			pos = findOffset(jsonContent.Bytes(), int(tErr.Offset))
		}
		if pos != nil {
			return nil, newError("failed to read config file at line ", pos.line, " char ", pos.char).Base(err)
		}
		return nil, newError("failed to read config file").Base(err)
	}

	pbConfig, err := jsonConfig.Build()
	if err != nil {
		return nil, newError("failed to parse json config").Base(err)
	}

	return pbConfig, nil
}
