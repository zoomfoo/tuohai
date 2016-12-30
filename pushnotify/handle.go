package pushnotify

import (
	"encoding/json"
)

type Handle struct {
	token, alert, sound, topic string
	badge                      int
}

func NewHandle(token string) *Handle {
	return &Handle{
		token: token,
	}
}

func (api *Handle) Aps() []byte {
	payload := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert": api.alert,
			"badge": api.badge,
			"sound": api.sound,
		},
	}

	b, _ := json.Marshal(payload)
	return b
}

func (api *Handle) Token() string {
	return api.token
}

func (api *Handle) Topic() string {
	// return api.topic
	return "com.yunzujia.woke"
}
