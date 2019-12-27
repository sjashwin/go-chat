package chat

// Message protocol passed over the network
type Message struct {
	ContentType string
	Data        string
	From, To    string
}
