package pkg

func RegisterAppConfigStore(store AppConfigStore) error {
	if appConfigStore != nil {
		return Error_RepeateRegisterStore
	}
	appConfigStore = store
	return nil
}
