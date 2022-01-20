package Vigor

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

var gauges = map[string]prometheus.Gauge{}
var statusGauge *prometheus.GaugeVec
var streamGaugeList = map[string]string{
	"actual":     "Downstream sync rate in bits per second, limited by ISP's line profile",
	"attainable": "maximum physical rate achievable in bits per second",
}
var endGaugeList = map[string]string{
	"snr_margin":     "SNR (Signal to Noise Ratio) value in dB",
	"attenuation":    "Line / Loop attenuation",
	"crc":            "Cyclic Redundancy Check fail count",
	"fecs":           "Forward Error Correction Seconds",
	"es":             "Erroed Seconds",
	"ses":            "Severely Erroed Seconds",
	"loss":           "Loss Of Signal Seconds",
	"uas":            "Un-Available Seconds",
	"hec_errors":     "Header Error Check Error count, HEC anomalies in the ATM Data Path",
	"rs_corrections": "RS Corrections - Not used",
	"los_failure":    "Loss of Signal Count",
	"lof_failure":    "Loss of Frame Count",
	"lpr_failure":    "Loss of Power Count",
	"ncd_failure":    "No Cell Delineation failure count",
	"lcs_failure":    "Loss Of Cell Delineation failure count",
	"nfec":           "Reed-Solomon codeword size in bytes used in the latency path in which the bearer channel is transported",
	"rfec":           "Actual number of Reed-Solomon redundancy bytes",
	"lysmb":          "Actual number of bits per symbol assigned to the latency path in which the bearer channel is transported",
}

func init() {
	for k, v := range streamGaugeList {
		name := fmt.Sprintf("vigor_%s_down", k)
		gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{Name: name, Help: v + " (Down)"})
		prometheus.MustRegister(gauges[name])
		name = fmt.Sprintf("vigor_%s_up", k)
		gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{Name: name, Help: v + " (Up)"})
		prometheus.MustRegister(gauges[name])
	}
	for k, v := range endGaugeList {
		name := fmt.Sprintf("vigor_%s_near", k)
		gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{Name: name, Help: v + " (Near)"})
		prometheus.MustRegister(gauges[name])
		name = fmt.Sprintf("vigor_%s_far", k)
		gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{Name: name, Help: v + " (Far)"})
		prometheus.MustRegister(gauges[name])
	}

	statusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vigor_status",
			Help: "vigor_status",
		},
		[]string{
			"firmware",
			"running_mode",
			"line_state",
			"power_mngt_mode",
			"vendor_id_modem",
			"vendor_id_dslam",
		})
	prometheus.MustRegister(statusGauge)
}
