package Vigor

import (
	"github.com/prometheus/client_golang/prometheus"
	"fmt"
)

var gauges = map[string]prometheus.Gauge{}
var statusGauge *prometheus.GaugeVec
var streamGaugeList = []string{
	"vigor_actual",
	"vigor_attainable",
}
var endGaugeList = []string{
	"vigor_snr_margin",
	"vigor_attenuation",
	"vigor_crc",
	"vigor_fecs",
	"vigor_es",
	"vigor_ses",
	"vigor_loss",
	"vigor_uas",
	"vigor_hec_errors",
	"vigor_rs_corrections",
	"vigor_los_failure",
	"vigor_lof_failure",
	"vigor_lpr_failure",
	"vigor_ncd_failure",
	"vigor_lcs_failure",
	"vigor_nfec",
	"vigor_rfec",
	"vigor_lysmb",
}

func init() {
	for _, prefix := range streamGaugeList {
		name := fmt.Sprintf("%s_down", prefix)
		gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{Name: name, Help: name})
		prometheus.MustRegister(gauges[name])
		name = fmt.Sprintf("%s_up", prefix)
		gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{Name: name, Help: name})
		prometheus.MustRegister(gauges[name])
	}
	for _, prefix := range endGaugeList {
		name := fmt.Sprintf("%s_near", prefix)
		gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{Name: name, Help: name})
		prometheus.MustRegister(gauges[name])
		name = fmt.Sprintf("%s_far", prefix)
		gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{Name: name, Help: name})
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
