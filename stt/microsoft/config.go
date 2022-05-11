package microsoft

import "encoding/json"

func ConfigFromJson(data []byte) Config {
	var c Config
	json.Unmarshal(data, &c)
	return c
}
