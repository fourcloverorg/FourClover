package common

import (
	"fmt"

	"github.com/fatih/color"
)

// App version
var (
	APP_VERSION     = "0.1.0"
	APP_NAME        = "fourclover"
	APP_DESCRIPTION = "OWASP Four Clover is a tool to help you test your applications and devices for security and policy issues."
	APP_AUTHOR      = "https://owasp.org/www-project-four-clover/"
	APP_URL         = "https://fourclover.org/"
)

// Encryption key
var (
	ci_key     = []byte("H0oJoO4hBvl")
	ph_key     = []byte("P58y9RR_zI")
	er_key     = []byte("BgyDrDrK7p4")
	CIPHER_KEY = []byte(fmt.Sprintf("%s%s%s", ci_key, ph_key, er_key))
)

// Colors
var (
	GREEN     = color.New(color.FgGreen).PrintlnFunc()
	YELLOW    = color.New(color.FgYellow).PrintlnFunc()
	MAGENTA   = color.New(color.FgMagenta).PrintlnFunc()
	RED       = color.New(color.FgRed).PrintlnFunc()
	INFO      = color.New(color.FgCyan).PrintlnFunc()
	INFO_BOLD = color.New(color.FgCyan, color.Bold).PrintlnFunc()
	WHITE     = color.New(color.FgWhite).PrintlnFunc()
)
