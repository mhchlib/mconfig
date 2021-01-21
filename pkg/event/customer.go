package event

import (
	"context"
	log "github.com/mhchlib/logger"
)

type EventCustomer struct {
	eventBus chan MConfigValEvent
}

func (e EventCustomer) AddEvent(event MConfigValEvent) error {
	e.eventBus <- event
	return nil
}

func (e EventCustomer) handleEvent(ctx context.Context) {
	log.Info("receive app config change event is started ")
	defer func() {
		log.Error("receive app config change event is closed ")
	}()
	for {
		select {
		case v, ok := <-e.eventBus:
			if !ok {
				return
			}
			log.Info("receive app ", v.appKey, " config change event ")
			//config 2 cache
			appConfigs := v.AppConfigs
			err := pkg.mconfigCache.putConfigCache(v.Key, appConfigs)
			if err != nil {
				log.Error(err)
				break
			}
			//notify client
			pkg.ConfigChangeChan <- v.Key
			log.Info("push app ", v.Key, " config change event to cache ")
		case <-ctx.Done():
			return
		}
	}
}
