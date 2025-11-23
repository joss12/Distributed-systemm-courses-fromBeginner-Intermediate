package pinger

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

func WriteCSV(w io.Writer, rs []Result) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := cw.Write([]string{"timestamp", "url", "status", "latency_ms", "error"}); err != nil {
		return err
	}
	for _, r := range rs {
		row := []string{
			r.Timestamps.Format("2006-01-02T15:04:05.000Z07:00"),
			r.URL,
			strconv.Itoa(r.Status),
			fmt.Sprintf("%.3f", float64(r.Latency.Milliseconds())/1000.0),
			r.Error,
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	return nil
}
