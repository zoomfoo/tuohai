package pushnotify

type Notification struct {
	Ios []Message `json:"ios"`
}

type Message struct {
	Token []string `json:"token"`
	Alert string   `json:"alert"`
	Sound string   `json:"sound"`
	Badge int      `json:"badge"`
}
