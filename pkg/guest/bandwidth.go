package guest

import (
	"fmt"
)

const Mb = 1000
const MB = Mb * 8
const Gb = Mb * 1000
const GB = Gb * 8
const Tb = Gb * 1000
const TB = Tb * 8

type Bandwidth struct {
	Value float64
	Unit  string
}

func (bandwidth Bandwidth) ToKbits() float64 {
	switch bandwidth.Unit {
	case "Mbits":
		return bandwidth.Value * Mb
	case "Gbits":
		return bandwidth.Value * Gb
	case "Tbits":
		return bandwidth.Value * Gb
	case "KBytes":
		return bandwidth.Value * 8
	case "MBytes":
		return bandwidth.Value * MB
	case "GBytes":
		return bandwidth.Value * GB
	case "TBytes":
		return bandwidth.Value * GB
	default:
		return bandwidth.Value
	}
}
func (bandwidth Bandwidth) String() string {
	return fmt.Sprintf("%.2f %s/sec", bandwidth.Value, bandwidth.Unit)
}

func HumanParseBandwidth(valueKb float64) string {
	switch {
	case valueKb >= Gb:
		return fmt.Sprintf("%.2f Gbits/sec", valueKb/Gb)
	case valueKb >= Mb:
		return fmt.Sprintf("%.2f Mbits/sec", valueKb/Mb)
	default:
		return fmt.Sprintf("%.2f Kbits/sec", valueKb)
	}
}
