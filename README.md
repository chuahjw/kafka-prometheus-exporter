# kafka-prometheus-exporter
Export Prometheus remote-write messages from Kafka to Prometheus or VictoriaMetrics

Based on: https://developer.confluent.io/get-started/go/

To be used with: https://github.com/chuahjw/prometheus-kafka-adapter

Build:
```
go mod init kafka-prometheus-exporter
go get github.com/confluentinc/confluent-kafka-go/kafka
CGO_ENABLED=1 go build -o out/kafka-prometheus-exporter util.go consumer.go
./out/kafka-prometheus-exporter kafka-prometheus-exporter.properties http://<TARGET>:<PORT>/api/v1/write
```