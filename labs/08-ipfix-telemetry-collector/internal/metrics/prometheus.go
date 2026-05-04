package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// PacketsReceived counts the total UDP IPFIX packets received
	PacketsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ipfix_udp_packets_received_total",
		Help: "The total number of UDP IPFIX packets received",
	})

	// FlowsProcessed counts the total IPFIX data records (flows) processed
	FlowsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ipfix_flows_processed_total",
		Help: "The total number of IPFIX data records processed",
	})

	// BytesProcessed counts the total bytes recorded in the flows, labeled by Src and Dst IP
	BytesProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ipfix_bytes_processed_total",
		Help: "Total bytes processed in IPFIX flows",
	}, []string{"src_ip", "dst_ip"})

	// ParseErrors counts decoding errors
	ParseErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ipfix_parse_errors_total",
		Help: "The total number of IPFIX parsing errors",
	})
)
