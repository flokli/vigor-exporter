package Vigor

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var ErrUpdateFailed = errors.New("dsl status update failed")
var ErrParseFailed = errors.New("dsl status parse failed")

var stream_labelmap = map[string]string{
	"Actual Rate":     "vigor_actual",
	"Attainable Rate": "vigor_attainable",
}

var end_labelmap = map[string]string{
	"SNR Margin":     "vigor_snr_margin",
	"Attenuation":    "vigor_attenuation",
	"CRC":            "vigor_crc",
	"FECS":           "vigor_fecs",
	"ES":             "vigor_es",
	"SES":            "vigor_ses",
	"LOSS":           "vigor_loss",
	"UAS":            "vigor_uas",
	"HEC Errors":     "vigor_hec_errors",
	"RS Corrections": "vigor_rs_corrections",
	"LOS Failure":    "vigor_los_failure",
	"LOF Failure":    "vigor_lof_failure",
	"LPR Failure":    "vigor_lpr_failure",
	"NCD Failure":    "vigor_ncd_failure",
	"LCD Failure":    "vigor_lcs_failure",
	"NFEC":           "vigor_nfec",
	"RFEC":           "vigor_rfec",
	"LYSMB":          "vigor_lysmb",
}

func (v *Vigor) UpdateStatus() error {
	resp, err := v.client.Get(fmt.Sprintf("http://%s/cgi-bin/V2X00.cgi?sFormAuthStr=%s&fid=2356", v.ip, v.csrf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 301 || resp.StatusCode == 302 {
		return nil
	}

	return ErrUpdateFailed
}

func (v *Vigor) FetchStatus() error {
	resp, err := v.client.Get(fmt.Sprintf("http://%s/doc/dslstatus.sht", v.ip))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ErrUpdateFailed
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return parseHTML(&bytes)
}

func stripHTML(input *[]byte) (*[]byte, error) {
	stripStuff := regexp.MustCompile(`(?m)(<[a-zA-Z]*) (.*?)>`)
	stripped := stripStuff.ReplaceAll(*input, []byte("$1>"))
	stripStuff = regexp.MustCompile(`(?m).*?<table>(.*)</table>.*`)
	stripped = stripStuff.ReplaceAll(stripped, []byte("$1"))
	stripStuff = regexp.MustCompile(`(?m)</?font>`)
	stripped = stripStuff.ReplaceAll(stripped, []byte(""))

	return &stripped, nil

}

func parsecol(input *[]byte, field string, multiplier float64, multicol bool) (float64, float64, error) {
	var colRegexp *regexp.Regexp
	if multicol {
		colRegexp = regexp.MustCompile(`<td>` + field + `</td><td>(.*?)</td><td>(.*?)</td>`)
	} else {

		colRegexp = regexp.MustCompile(`<td>` + field + `</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td>`)
	}

	col := map[int]float64{}
	if multiplier == 0 {
		multiplier = 1
	}

	for i, match := range colRegexp.FindSubmatch(*input) {
		if i == 0 {
			continue
		}

		value, _ := strconv.ParseFloat(strings.Replace(string(match), " ", "", -1), 32)
		col[i-1] = value * multiplier
	}

	if len(col) != 2 {
		return 0, 0, ErrParseFailed
	}

	return col[0], col[1], nil
}

func parseHeadCol(input *[]byte, field string) (string, string, error) {
	colRegexp := regexp.MustCompile(`<td>` + field + `:</td><td>(.*?)</td>`)

	value1 := "n/a"
	value2 := "n/a"
	for i, match := range colRegexp.FindAllSubmatch(*input, 2) {
		for i2, match2 := range match {
			switch i2 {
			case 0:
				continue
			case 1:
				if i == 0 {
					value1 = strings.Replace(string(match2), "&nbsp;", " ", -1)
				} else {
					value2 = strings.Replace(string(match2), "&nbsp;", " ", -1)
				}
			}
		}
	}
	return value1, value2, nil
}

func parseHTML(input *[]byte) error {

	stripped, err := stripHTML(input)
	if err != nil {
		return err
	}

	for htmlField, gaugePrefix := range stream_labelmap {
		var multicol bool = false
		var multiplier float64 = 1
		switch htmlField {
		case "Attainable Rate",
			"Actual Rate":
			multiplier = 1000
		case
			"Interleave Depth":
			multicol = true
		}

		down, up, err := parsecol(stripped, htmlField, multiplier, multicol)
		if err != nil {
			return err
		}

		gauges[gaugePrefix+"_down"].Set(down)
		gauges[gaugePrefix+"_up"].Set(up)
	}

	for htmlField, gaugePrefix := range end_labelmap {
		var multicol bool = false
		var multiplier float64 = 1
		switch htmlField {
		case
			"CRC",
			"HEC Errors",
			"RS Corrections",
			"LOS Failure",
			"LOF Failure",
			"LPR Failure",
			"NCD Failure",
			"LCD Failure",
			"NFEC",
			"RFEC",
			"LYSMB":
			multicol = true
		}

		nearend, farend, err := parsecol(stripped, htmlField, multiplier, multicol)
		if err != nil {
			return err
		}

		gauges[gaugePrefix+"_near"].Set(nearend)
		gauges[gaugePrefix+"_far"].Set(farend)
	}

	firmware, _, _ := parseHeadCol(stripped, "Firmware")
	running_mode, _, _ := parseHeadCol(stripped, "Running Mode")
	line_state, _, _ := parseHeadCol(stripped, "Line State")
	power_mngt_mode, _, _ := parseHeadCol(stripped, "Power Mngt Mode")
	vendor_id_modem, vendor_id_dslam, _ := parseHeadCol(stripped, "Vendor ID")

	statusGauge.Reset()
	statusGauge.With(prometheus.Labels{
		"firmware":        firmware,
		"running_mode":    running_mode,
		"line_state":      line_state,
		"power_mngt_mode": power_mngt_mode,
		"vendor_id_modem": vendor_id_modem,
		"vendor_id_dslam": vendor_id_dslam,
	}).Set(1)

	return nil
}
