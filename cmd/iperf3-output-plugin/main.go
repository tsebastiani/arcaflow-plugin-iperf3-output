package main

import (
	arcaflow_plugin_iperf3_output "arcaflow-plugin-iperf3-output"
	"go.flow.arcalot.io/pluginsdk/plugin"
)

func main() {
	plugin.Run(arcaflow_plugin_iperf3_output.Schema)
}
