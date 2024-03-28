package main

type Memory struct {
	data map[string][]Transaction
}

func NewMemory() *Memory {
	return &Memory{data: make(map[string][]Transaction)}
}

func (m *Memory) AddSubscriber(address string) {
	m.data[address] = []Transaction{}
}

func (m *Memory) GetAllSubscribers() []string {
	subscribers := make([]string, 0, len(m.data))

	for address := range m.data {
		subscribers = append(subscribers, address)
	}

	return subscribers
}

func (m *Memory) SubscriberExist(address string) bool {
	_, ok := m.data[address]
	return ok
}

func (m *Memory) AddTransaction(address string, transaction Transaction) {
	m.data[address] = append(m.data[address], transaction)
}

func (m *Memory) GetSubscriberTransactions(address string) []Transaction {
	return m.data[address]
}
