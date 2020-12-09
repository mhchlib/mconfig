package pkg

import "errors"

var (
	// Error_RepeateRegisterStore ...
	Error_RepeateRegisterStore = errors.New("RepeateRegisterStore")
	// Error_AppConfigNotFound ...
	Error_AppConfigNotFound = errors.New("AppConfigNotFound")
	// Error_AppConfigByFilterNotFound ...
	Error_AppConfigByFilterNotFound = errors.New("AppConfigByFilterNotFound")
	// Error_ParserAppConfigFail ...
	Error_ParserAppConfigFail = errors.New("ParserAppConfigFail")
)
