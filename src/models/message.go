package models

type Message struct {
	ID        string `json:"id"`
	Source    string `json:"source"`
	Recipient string `json:"recipient"`
	Name      string `json:"name"`
	Timestamp uint64 `json:"timestamp"`
	Body      string `json:"body"`
}

type ByTimestamp []Message

func (m ByTimestamp) Len() int           { return len(m) }
func (m ByTimestamp) Less(i, j int) bool { return m[i].Timestamp > m[j].Timestamp }
func (m ByTimestamp) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
