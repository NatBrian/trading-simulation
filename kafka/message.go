package kafka

type Message struct {
	Value     []byte
	Headers   map[string][]byte
	Partition int
	Offset    int64
	Key       []byte
}
