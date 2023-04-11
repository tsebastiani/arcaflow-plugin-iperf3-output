package arcaflow_plugin_iperf3_output

import (
	"arcaflow-plugin-iperf3-output/pkg/iperf"
	"arcaflow-plugin-iperf3-output/pkg/result"
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"go.flow.arcalot.io/pluginsdk/schema"
	"io/ioutil"
	"os"
	"time"
)

var iperf3_schema = schema.NewScopeSchema(
	schema.NewStructMappedObjectSchema[IperfOutput](
		"IperfOutput",
		map[string]*schema.PropertySchema{
			"start": schema.NewPropertySchema(
				schema.NewAnySchema(),
				schema.NewDisplayValue(
					nil,
					nil,
					nil,
				),
				false,
				nil,
				nil,
				nil,
				schema.PointerTo(`"Streams"`),
				nil,
			).TreatEmptyAsDefaultValue(),
			"intervals": schema.NewPropertySchema(
				schema.NewAnySchema(),
				schema.NewDisplayValue(
					nil,
					nil,
					nil,
				),
				false,
				nil,
				nil,
				nil,
				schema.PointerTo(`"Streams"`),
				nil,
			).TreatEmptyAsDefaultValue(),
			"end": schema.NewPropertySchema(
				end_schema,
				schema.NewDisplayValue(
					nil,
					nil,
					nil,
				),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
)

var end_schema = schema.NewStructMappedObjectSchema[End](
	"End",
	map[string]*schema.PropertySchema{
		"streams": schema.NewPropertySchema(
			schema.NewAnySchema(),
			schema.NewDisplayValue(
				nil,
				nil,
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"Streams"`),
			nil,
		).TreatEmptyAsDefaultValue(),
		"sum_sent": schema.NewPropertySchema(
			schema.NewAnySchema(),
			schema.NewDisplayValue(
				nil,
				nil,
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"SumSent"`),
			nil,
		).TreatEmptyAsDefaultValue(),
		"sum_received": schema.NewPropertySchema(
			sum_received_schema,
			schema.NewDisplayValue(
				nil,
				nil,
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
		"cpu_utilization_percent": schema.NewPropertySchema(
			schema.NewAnySchema(),
			schema.NewDisplayValue(
				nil,
				nil,
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"CpuUtilizationPercent"`),
			nil,
		).TreatEmptyAsDefaultValue(),
		"sender_tcp_congestion": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				nil,
				nil,
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"receiver_tcp_congestion": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				nil,
				nil,
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
	},
)

var sum_received_schema = schema.NewStructMappedObjectSchema[SumReceived](
	"SumReceived",
	map[string]*schema.PropertySchema{
		"start": schema.NewPropertySchema(
			schema.NewFloatSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Start"),
				nil,
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"end": schema.NewPropertySchema(
			schema.NewFloatSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("End"),
				nil,
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"seconds": schema.NewPropertySchema(
			schema.NewFloatSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Seconds"),
				nil,
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"bytes": schema.NewPropertySchema(
			schema.NewFloatSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Bytes"),
				nil,
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"bits_per_second": schema.NewPropertySchema(
			schema.NewFloatSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("BitsPerSecond"),
				nil,
				nil,
			),
			true,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
	},
)

var Schema = schema.NewCallableSchema(
	schema.NewCallableStep[IperfOutput](
		"render",
		iperf3_schema,
		map[string]*schema.StepOutputSchema{
			"success": schema.NewStepOutputSchema(
				successOutputSchema,
				schema.NewDisplayValue(
					nil, nil,
					nil,
				),
				false,
			),
			"error": schema.NewStepOutputSchema(
				errorOutputSchema,
				schema.NewDisplayValue(
					nil, nil,
					nil,
				),
				true,
			),
		},
		schema.NewDisplayValue(
			nil, nil,
			nil,
		),
		func(input IperfOutput) (string, any) {
			json_new := jsoniter.ConfigCompatibleWithStandardLibrary
			var data []byte
			data, err := json_new.Marshal(input)

			var sr result.ScenarioResults
			npr := result.Data{}
			npr.StartTime = time.Now()
			npr.Driver = "iperf3"
			npr.Profile = "STREAM"

			if err != nil {
				return "error", ErrorOutput{
					err.Error(),
				}
			}
			buffer := bytes.NewBuffer(data)
			nr, err := iperf.ParseResults(buffer)
			if err != nil {
				return "error", ErrorOutput{
					err.Error(),
				}
			}

			npr.Metric = nr.Metric
			npr.Samples = 1
			npr.ThroughputSummary = append(npr.ThroughputSummary, nr.Throughput)
			npr.LatencySummary = append(npr.LatencySummary, nr.Latency99ptile)
			sr.Results = append(sr.Results, npr)

			//redirect stdout to string
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			result.ShowStreamResult(sr)
			//result.ShowRRResult(sr)
			//result.ShowLatencyResult(sr)

			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout

			outStr := string(out)
			npr.EndTime = time.Now()
			return "success", SuccessOutput{
				outStr,
			}
		},
	),
)
