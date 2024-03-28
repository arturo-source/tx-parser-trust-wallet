package main

// Memory represents an imaginary database
type Memory struct {
	data map[string][]Transaction
}

// NewMemory creates a new memory
func NewMemory() *Memory {
	return &Memory{data: make(map[string][]Transaction)}
}

// AddSubscriber adds a new subscriber to the memory
func (m *Memory) AddSubscriber(address string) {
	m.data[address] = []Transaction{}
}

// GetAllSubscribers returns all subscribers in the memory
func (m *Memory) GetAllSubscribers() []string {
	subscribers := make([]string, 0, len(m.data))

	for address := range m.data {
		subscribers = append(subscribers, address)
	}

	return subscribers
}

// SubscriberExist checks if an address is already in the memory
func (m *Memory) SubscriberExist(address string) bool {
	_, ok := m.data[address]
	return ok
}

// AddTransaction adds a new transaction to the memory with the given address
func (m *Memory) AddTransaction(address string, transaction Transaction) {
	m.data[address] = append(m.data[address], transaction)
}

// GetSubscriberTransactions returns all transactions for the given address
func (m *Memory) GetSubscriberTransactions(address string) []Transaction {
	return m.data[address]
}
