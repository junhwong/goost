package named

import "time"

var functions = map[string]func() (interface{}, error){
	"time_now": func() (interface{}, error) {
		return time.Now(), nil
	},
	"time_now_nano": func() (interface{}, error) {
		return time.Now().UnixNano(), nil
	},
	"timestamp_us": func() (interface{}, error) {
		return time.Now().UnixMicro(), nil
	},
}
