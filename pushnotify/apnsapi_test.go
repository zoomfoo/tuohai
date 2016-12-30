package pushnotify

import (
	"encoding/json"
	"testing"
)

func TestAPNs(t *testing.T) {
	js := `{
    "ios":[
        {
            	"token":"eac57ab70a9eaa8ff26a94a3057a2c639ae29bb213c824bde520d36e76f7c446",
            	"alert":"wangyang",
            	"badge":1,
            	"sound":"default"
        	}
    	]
	}`

	var n Notification
	err := json.Unmarshal([]byte(js), &n)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(n)
}

func TestUnmarshal(t *testing.T) {
	msg := Message{
		Token: []string{"eac57ab70a9eaa8ff26a94a3057a2c639ae29bb213c824bde520d36e76f7c446"},
		Alert: "2323",
		Badge: 1,
		Sound: "default",
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(b))
}
