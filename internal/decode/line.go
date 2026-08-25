package decode

import (
	"bufio"
	"bytes"

	"github.com/dyl-06/telemetry/internal/model"
)

// DecodeLines 解析 NDJSON 格式的多条点位。
func DecodeLines(data []byte) ([]model.Point, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	points := make([]model.Point, 0)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		point, err := Decode(line)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return points, nil
}
