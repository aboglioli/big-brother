package events

type Event struct {
	Type     string `json:"type"`
	Exchange string `json:"exchange"`
	Route    string `json:"route"`
	Queue    string `json:"queue"`
}
