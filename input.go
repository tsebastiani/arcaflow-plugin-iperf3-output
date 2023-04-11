package arcaflow_plugin_iperf3_output

type IperfOutput struct {
	Start     any `json:"start"`
	Intervals any `json:"intervals"`
	End       End `json:"end"`
}

type End struct {
	Streams               any         `json:"streams"`
	SumSent               any         `json:"sum_sent"`
	SumReceived           SumReceived `json:"sum_received"`
	CpuUtilizationPercent any         `json:"cpu_utilization_percent"`
	SenderTcpCongestion   string      `json:"sender_tcp_congestion"`
	ReceiverTcpCongestion string      `json:"receiver_tcp_congestion"`
}

type SumReceived struct {
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         float64 `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
}
