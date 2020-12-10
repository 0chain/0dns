module zdns

require (
	0dns.io/core v0.0.0
	0dns.io/zdnscore v0.0.0
	github.com/0chain/gosdk v1.1.8
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.15.0
)

replace 0dns.io/core => ../../core

replace 0dns.io/zdnscore => ../../zdnscore

go 1.14
