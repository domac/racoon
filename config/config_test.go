package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_load_config_by_string(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	conf, err := LoadConfigFromString(`{
		"input": [{
			"type": "exec",
			"command": "uptime",
			"interval": 3
		},{
			"type": "exec",
			"command": "whoami",
			"interval": 4
		}],
		"output": [{
			"type": "stdout",
            "host": "127.0.0.1"
		}]
	}`)

	require.NoError(err)
	require.NotNil(conf)

}
