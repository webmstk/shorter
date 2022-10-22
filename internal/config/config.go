package config

import "fmt"

const ServerPort = "8080"
const ServerHost = "localhost"
const ServerSchema = "http"

var ServerBaseURL = fmt.Sprintf("%s://%s:%s", ServerSchema, ServerHost, ServerPort)
