package microsoft

import "encoding/json"

func ConfigFromJson(id int, cb string, data []byte) Config {
	var c Config
	json.Unmarshal(data, &c)
	c.Callback = cb
	c.Id = id
	return c
}
