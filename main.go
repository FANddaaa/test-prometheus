package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus相关
const (
	Namespace = "g3Network"
	SubSystem = "nlb_monitor"
)
const (
	LabelBusID = "busID"
)
const (
	LabelType = "type"
	TypeConn  = "conn"
	TypePkt   = "pkt"
	TypeHc    = "hc"
)
const (
	LabelUnit      = "unit"
	UnitNewConn    = "new_conn"    //
	UnitCurrConn   = "curr_conn"   //
	UnitFailedConn = "failed_conn" //
	UnitBit        = "bit"         // 流量
	UnitBps        = "bps"         // 带宽
	UnitPkt        = "pkt"         // 包量
	UnitPps        = "pps"         // 每秒网络数据包数量
	UnitUnhealthy  = "unhealthy"   // 异常数
	UnitHealthy    = "healthy"     // 正常数
)
const (
	LabelDirection = "direction"
	DirectionIn    = "in"
	DirectionOut   = "out"
)

// 抓取信息时间间隔
const (
	interval = 5 // unit: second
)

var (
	NlbMessage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: SubSystem,
			Name:      "nlb_message",
		},
		[]string{LabelBusID, LabelType, LabelUnit, LabelDirection},
	)
)

func main() {
	tick := time.NewTicker(time.Second * interval)
	defer tick.Stop()
	prometheus.MustRegister(NlbMessage)
	go upload()
	busID := "nlb-xxxxxx"

	for {
		select {
		case <-tick.C:
			InPkts := rand.Int() % 1000
			OutPkts := rand.Int() % 1000
			InBits := rand.Int() % 1000
			OutBits := rand.Int() % 1000
			InBps := InBits / interval
			OutBps := OutBits / interval
			NlbMessage.WithLabelValues(busID, TypePkt, UnitPkt, DirectionIn).Add(float64(InPkts))
			NlbMessage.WithLabelValues(busID, TypePkt, UnitPkt, DirectionOut).Add(float64(OutPkts))
			NlbMessage.WithLabelValues(busID, TypePkt, UnitBit, DirectionIn).Add(float64(InBits))
			NlbMessage.WithLabelValues(busID, TypePkt, UnitBit, DirectionOut).Add(float64(OutBits))
			NlbMessage.WithLabelValues(busID, TypePkt, UnitBps, DirectionIn).Set(float64(InBps))
			NlbMessage.WithLabelValues(busID, TypePkt, UnitBps, DirectionOut).Set(float64(OutBps))
		}
	}
}

func upload() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
