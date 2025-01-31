package network

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Metric used in monitoring service.
var (
	estimatedNetworkSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Help:      "Estimated network size",
			Name:      "network_size",
			Namespace: "neogo",
		},
	)

	peersConnected = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Help:      "Number of connected peers",
			Name:      "peers_connected",
			Namespace: "neogo",
		},
	)

	servAndNodeVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Help:      "Server and Node versions",
			Name:      "serv_node_version",
			Namespace: "neogo",
		},
		[]string{"description", "value"},
	)

	poolCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Help:      "Number of available node addresses",
			Name:      "pool_count",
			Namespace: "neogo",
		},
	)

	blockQueueLength = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Help:      "Block queue length",
			Name:      "block_queue_length",
			Namespace: "neogo",
		},
	)
	p2pCmds = make(map[CommandType]prometheus.Histogram)
)

func init() {
	prometheus.MustRegister(
		estimatedNetworkSize,
		peersConnected,
		servAndNodeVersion,
		poolCount,
		blockQueueLength,
	)
	for _, cmd := range []CommandType{CMDVersion, CMDVerack, CMDGetAddr,
		CMDAddr, CMDPing, CMDPong, CMDGetHeaders, CMDHeaders, CMDGetBlocks,
		CMDMempool, CMDInv, CMDGetData, CMDGetBlockByIndex, CMDNotFound,
		CMDTX, CMDBlock, CMDExtensible, CMDP2PNotaryRequest, CMDGetMPTData,
		CMDMPTData, CMDReject, CMDFilterLoad, CMDFilterAdd, CMDFilterClear,
		CMDMerkleBlock, CMDAlert} {
		p2pCmds[cmd] = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Help:      "P2P " + cmd.String() + " handling time",
				Name:      "p2p_" + strings.ToLower(cmd.String()) + "_time",
				Namespace: "neogo",
			},
		)
		prometheus.MustRegister(p2pCmds[cmd])
	}
}

func updateNetworkSizeMetric(sz int) {
	estimatedNetworkSize.Set(float64(sz))
}

func updateBlockQueueLenMetric(bqLen int) {
	blockQueueLength.Set(float64(bqLen))
}

func updatePoolCountMetric(pCount int) {
	poolCount.Set(float64(pCount))
}

func updatePeersConnectedMetric(pConnected int) {
	peersConnected.Set(float64(pConnected))
}
func setServerAndNodeVersions(nodeVer string, serverID string) {
	servAndNodeVersion.WithLabelValues("Node version: ", nodeVer).Add(0)
	servAndNodeVersion.WithLabelValues("Server id: ", serverID).Add(0)
}
func addCmdTimeMetric(cmd CommandType, t time.Duration) {
	// Shouldn't happen, message decoder checks the type, but better safe than sorry.
	if p2pCmds[cmd] == nil {
		return
	}
	p2pCmds[cmd].Observe(t.Seconds())
}
