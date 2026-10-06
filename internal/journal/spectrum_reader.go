package journal

import (
	"bufio"
	"encoding/json"
	"io"
)

type SpectrumRecord struct {
	Wavelength []float64 `json:"wavelength"`
	Intensity  []float64 `json:"intensity"`
}

func ReadSpectrum(r io.Reader) (SpectrumRecord, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 4*1024*1024)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return SpectrumRecord{}, err
		}
		return SpectrumRecord{}, io.EOF
	}
	var record SpectrumRecord
	if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
		return SpectrumRecord{}, err
	}
	if err := scanner.Err(); err != nil {
		return SpectrumRecord{}, err
	}
	return record, nil
}
