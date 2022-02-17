package infrastructure

type AggregateTests struct {
	dispatcher  Dispatcher
	store       *FakeAggregateStore
	latestError error
}

func NewAggregateTests(store *FakeAggregateStore) AggregateTests {
	return AggregateTests{
		store: store,
	}
}

func (t *AggregateTests) RegisterHandlers(handlers CommandHandler) {
	commandHandlerMap := NewCommandHandlerMap(handlers)
	t.dispatcher = NewDispatcher(commandHandlerMap)
}

func (t *AggregateTests) Given(events... interface{}) {
	t.store.SetInitialEvents(events)
}

func (t *AggregateTests) When(command interface{}) {
	t.latestError = t.dispatcher.Dispatch(command)
}

func (t *AggregateTests) Then(then func([]interface{}, error)) {
	then(t.store.GetStoredChanges(), t.latestError)
}
