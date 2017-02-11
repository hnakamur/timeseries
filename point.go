package timeseries

type Point struct {
	Timestamp uint32
	Value     float64
}

//func Marshal(points []Point) ([]byte, error) {
//	var buf bytes.Buffer
//	return nil, nil
//}
//
//func Unmarshal(data []byte) ([]Point, error) {
//	return nil, nil
//}
