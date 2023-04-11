package arcaflow_plugin_iperf3_output

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSchemaValidation(t *testing.T) {
	data, err := os.ReadFile("testdata/schema.json")
	assert.Nil(t, err)
	var input = IperfOutput{}
	json.Unmarshal(data, &input)
	fmt.Println(input)
	sc, err := Schema.SelfSerialize()
	assert.Nil(t, err)
	err = iperf3_schema.Validate(input)
	assert.Nil(t, err)
	fmt.Println(sc)

	serialized_schema, err := iperf3_schema.Serialize(input)
	assert.Nil(t, err)
	outputId, serializedOutput, err := Schema.Call("render", serialized_schema)

	assert.Nil(t, err)
	fmt.Println(outputId)
	fmt.Println(serializedOutput)
}
