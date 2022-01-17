package yeelight

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Parse_Discovery(t *testing.T) {

	sampleResp :=
		`HTTP/1.1 200 OK
		Cache-Control: max-age=3600
		Date: 
		Ext: 
		Location: yeelight://192.168.15.58:55443
		Server: POSIX UPnP/1.0 YGLC/1
		id: 0x0000000007fe356e
		model: color
		fw_ver: 65
		support: get_prop set_default set_power toggle set_bright start_cf stop_cf set_scene cron_add cron_get cron_del set_ct_abx set_rgb set_hsv set_adjust adjust_bright adjust_ct adjust_color set_music set_name
		power: on
		bright: 100
		color_mode: 3
		ct: 5004
		rgb: 65348
		hue: 136
		sat: 100
		name: 
		`
	expectedYeelight := Yeelight{
		id:      "0x0000000007fe356e",
		name:    "0x0000000007fe356e",
		addr:    "192.168.15.58:55443",
		model:   "color",
		support: []string{"get_prop", "set_default", "set_power", "toggle", "set_bright", "start_cf", "stop_cf", "set_scene", "cron_add", "cron_get", "cron_del", "set_ct_abx", "set_rgb", "set_hsv", "set_adjust", "adjust_bright", "adjust_ct", "adjust_color", "set_music", "set_name"},
	}

	resultYeelight := parseDiscoveyResponse(sampleResp)

	assert.Equal(t, expectedYeelight, resultYeelight)
}
