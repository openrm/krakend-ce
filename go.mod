module github.com/devopsfaith/krakend-ce

go 1.12

require (
	github.com/devopsfaith/bloomfilter v1.4.0
	github.com/devopsfaith/krakend v1.2.0
	github.com/devopsfaith/krakend-amqp v1.4.0
	github.com/devopsfaith/krakend-botdetector v1.4.0
	github.com/devopsfaith/krakend-cel v1.4.0
	github.com/devopsfaith/krakend-circuitbreaker v1.4.0
	github.com/devopsfaith/krakend-cobra v1.4.0
	github.com/devopsfaith/krakend-consul v1.4.0
	github.com/devopsfaith/krakend-cors v1.4.0
	github.com/devopsfaith/krakend-etcd v0.0.0-20190425091451-d989a26508d7
	github.com/devopsfaith/krakend-flexibleconfig v1.4.0
	github.com/devopsfaith/krakend-gelf v1.4.0
	github.com/devopsfaith/krakend-gologging v1.4.0
	github.com/devopsfaith/krakend-httpcache v1.4.0
	github.com/devopsfaith/krakend-httpsecure v1.4.0
	github.com/devopsfaith/krakend-influx v1.4.0
	github.com/devopsfaith/krakend-jose v1.4.0
	github.com/devopsfaith/krakend-jsonschema v1.4.0
	github.com/devopsfaith/krakend-lambda v1.4.0
	github.com/devopsfaith/krakend-logstash v1.4.0
	github.com/devopsfaith/krakend-lua v1.4.0
	github.com/devopsfaith/krakend-martian v1.4.0
	github.com/devopsfaith/krakend-metrics v1.4.0
	github.com/devopsfaith/krakend-oauth2-clientcredentials v1.4.0
	github.com/devopsfaith/krakend-opencensus v1.4.1
	github.com/devopsfaith/krakend-pubsub v1.4.0
	github.com/devopsfaith/krakend-ratelimit v1.4.0
	github.com/devopsfaith/krakend-rss v1.4.0
	github.com/devopsfaith/krakend-usage v1.4.0
	github.com/devopsfaith/krakend-viper v1.4.0
	github.com/devopsfaith/krakend-xml v1.4.0
	github.com/gin-gonic/gin v1.7.2
	github.com/go-contrib/uuid v1.2.0
	github.com/hashicorp/vault v1.6.0 // indirect
	github.com/influxdata/influxdb v1.7.4 // indirect
	github.com/kpacha/opencensus-influxdb v0.0.0-20181102202715-663e2683a27c // indirect
	github.com/letgoapp/krakend-consul v0.0.0-20190130102841-7623a4da32a1 // indirect
	github.com/luraproject/lura v1.4.1
	github.com/openrm/krakend-bloomd v0.0.0-20210716064324-815ce3f30e0a
	github.com/openrm/krakend-sentry v0.0.0-20210719091511-65b801de7883
	github.com/openrm/module-tracing-golang v1.1.2
	gocloud.dev/pubsub/kafkapubsub v0.21.0 // indirect
	gocloud.dev/pubsub/natspubsub v0.21.0 // indirect
	gocloud.dev/pubsub/rabbitpubsub v0.21.0 // indirect
	gocloud.dev/secrets/hashivault v0.21.0 // indirect
	k8s.io/api v0.20.2 // indirect
)

replace github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 => github.com/m4ns0ur/httpcache v0.0.0-20200426190423-1040e2e8823f

replace (
	github.com/devopsfaith/krakend-amqp v1.4.0 => github.com/openrm/krakend-amqp v0.0.0-20210929014739-e4cde5f13338
	github.com/devopsfaith/krakend-cel v1.4.0 => github.com/openrm/krakend-cel v0.0.0-20210719101150-9ec8a72804f1
	github.com/devopsfaith/krakend-jose v1.4.0 => github.com/openrm/krakend-jose v0.0.0-20210719112350-52718b067cbd
	github.com/devopsfaith/krakend-martian v1.4.0 => github.com/openrm/krakend-martian v0.0.0-20210718024057-b5801065122d
	github.com/devopsfaith/krakend-opencensus v1.4.1 => github.com/openrm/krakend-opencensus v0.0.0-20210831072122-582d5d334306
)
