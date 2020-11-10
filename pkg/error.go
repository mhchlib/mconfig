package pkg

import "errors"

var (
	Error_RepeateRegisterStore      = errors.New("RepeateRegisterStore")
	Error_AppConfigNotFound         = errors.New("AppConfigNotFound")
	Error_AppConfigByFilterNotFound = errors.New("AppConfigByFilterNotFound")
	Error_ParserAppConfigFail       = errors.New("ParserAppConfigFail")
)
