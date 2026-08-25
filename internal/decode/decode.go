package decode

import (
	"encoding/json"
	"errors"

	"github.com/dyl-06/telemetry/internal/model"
)

type wirePoint struct {
	Series string  `json:"series"`
	Seq    int64   `json:"seq"`
	Value  float64 `json:"value"`
	Ts     int64   `json:"ts"`
}

// Decode 解析一个点位；字段类型异常时返回错误。
func Decode(raw []byte) (model.Point, error) {
	var wp wirePoint
	if err := json.Unmarshal(raw, &wp); err != nil {
		return model.Point{}, err
	}
	if wp.Series == "" {
		return model.Point{}, errors.New("missing series")
	}
	return model.Point{
		Series: wp.Series,
		Seq:    wp.Seq,
		Value:  wp.Value,
		Ts:     wp.Ts,
	}, nil
}
