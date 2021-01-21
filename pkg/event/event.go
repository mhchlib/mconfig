package event

import (
	"context"
)

func NewMConfigEventCustomer() MConfigEventCustomer {
	customer = &EventCustomer{
		eventBus: make(chan MConfigValEvent, LENGTH_MAX_EVENT),
	}
	return customer
}

func StartMConfigStoreEventBus(ctx context.Context) {
	customer.handleEvent(ctx)
}
