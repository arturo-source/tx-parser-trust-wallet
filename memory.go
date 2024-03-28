package main

type Memory map[string][]Transaction

func NewMemory() Memory {
	return make(Memory)
}

func (m Memory) AddSubscriber(address string) {
	m[address] = []Transaction{}
}

func (m Memory) GetAllSubscribers() []string {
	subscribers := make([]string, 0, len(m))

	for address := range m {
		subscribers = append(subscribers, address)
	}

	return subscribers
}

func (m Memory) SubscriberExist(address string) bool {
	_, ok := m[address]
	return ok
}

func (m Memory) AddTransaction(address string, transaction Transaction) {
	m[address] = append(m[address], transaction)
}

func (m Memory) GetSubscriberTransactions(address string) []Transaction {
	return m[address]
}
