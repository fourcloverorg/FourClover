package commands

import (
	"log"

	common "fourclover.org/internal/common"
)

func Version() bool {
	log.Default().Println("INFO: FourClover version is :", common.APP_VERSION)
	return true
}
