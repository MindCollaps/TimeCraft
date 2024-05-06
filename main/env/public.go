package env

import "embed"

var UNIX = false
var DEBUG = false
var TESTING = false

var Files embed.FS

const VERSION = "v0.0.1"
const BANNER = ""
