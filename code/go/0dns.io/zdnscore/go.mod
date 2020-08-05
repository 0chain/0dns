module 0dns.io/zdnscore

require (
	0dns.io/core v0.0.0
	github.com/0chain/gosdk v1.0.94
	go.mongodb.org/mongo-driver v1.4.0
	go.uber.org/zap v1.15.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

replace 0dns.io/core => ../core

go 1.14
