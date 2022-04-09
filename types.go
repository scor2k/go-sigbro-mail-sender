package main

// Template for every email via RabbitMQ
type emailJSON struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Sender    string `json:"sender"`
}
