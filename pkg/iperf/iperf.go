package iperf

import (
	"bytes"
	"encoding/json"

	log "arcaflow-plugin-iperf3-output/pkg/logging"
	"arcaflow-plugin-iperf3-output/pkg/sample"
)

const workload = "iperf3"

type Result struct {
	Data struct {
		TCPStream struct {
			Rate float32 `json:"bits_per_second"`
		} `json:"sum_received"`
		UDPStream struct {
			Rate        float32 `json:"bits_per_second"`
			LossPercent float32 `json:"lost_percent"`
		} `json:"sum"`
	} `json:"end"`
}

// ParseResults accepts the stdout from the execution of the benchmark. It also needs
// The NetConfig to determine aspects of the workload the user provided.
// It will return a Sample struct or error
func ParseResults(stdout *bytes.Buffer) (sample.Sample, error) {
	sample := sample.Sample{}
	sample.Driver = workload
	result := Result{}
	sample.Metric = "Mb/s"
	json.NewDecoder(stdout).Decode(&result)
	if result.Data.TCPStream.Rate > 0 {
		sample.Throughput = float64(result.Data.TCPStream.Rate) / 1000000

	} else {
		sample.Throughput = float64(result.Data.UDPStream.Rate) / 1000000
	}
	log.Debugf("Storing %s sample throughput:  %f", sample.Driver, sample.Throughput)

	return sample, nil
}
