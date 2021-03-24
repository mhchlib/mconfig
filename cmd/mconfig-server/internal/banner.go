package internal

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/version"
)

var banner = `
     _____ ______    ________   ________   ________    ________  ___   ________
    |\   _ \  _   \ |\   ____\ |\   __  \ |\   ___  \ |\  _____\|\  \ |\   ____\
    \ \  \\\__\ \  \\ \  \___| \ \  \|\  \\ \  \\ \  \\ \  \__/ \ \  \\ \  \___|
     \ \  \\|__| \  \\ \  \     \ \  \\\  \\ \  \\ \  \\ \   __\ \ \  \\ \  \  ___
      \ \  \    \ \  \\ \  \____ \ \  \\\  \\ \  \\ \  \\ \  \_|  \ \  \\ \  \|\  \
       \ \__\    \ \__\\ \_______\\ \_______\\ \__\\ \__\\ \__\    \ \__\\ \_______\
        \|__|     \|__| \|_______| \|_______| \|__| \|__| \|__|     \|__| \|_______|

`

// ShowBanner ...
func ShowBanner() {
	banner = banner + "mconfig-server made by QMS \n"
	banner = banner + "version: " + version.GetVersion() + "\n"
	log.Info(banner)
}
