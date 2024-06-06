package env

import (
	"embed"
	"time"
)

var UNIX = false
var DEBUG = false
var TESTING = false

var Files embed.FS

var TimeZone *time.Location

const VERSION = "v0.0.1"
const BANNER = ""
