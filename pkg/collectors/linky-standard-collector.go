package collectors

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

// LinkyCollector object to describe and collect metrics
type LinkyCollector struct {
	device                 string
	baudRate               int
	frameSize              byte
	parity                 serial.Parity
	stopBits               serial.StopBits
	linky_info             *prometheus.Desc
	linky_index            *prometheus.Desc
	linky_current          *prometheus.Desc
	linky_voltage          *prometheus.Desc
	linky_subscribed_power *prometheus.Desc
	linky_power            *prometheus.Desc
	linky_load_management  *prometheus.Desc
	linky_relays           *prometheus.Desc
	linky_provider_day     *prometheus.Desc
}

// Internal linky values object to each metrics
type linkyValues struct {
	adsc                string
	vtic                string
	date                string
	ngtf                string
	ltarf               string
	east                uint32
	easf01              uint32
	easf02              uint32
	easf03              uint32
	easf04              uint32
	easf05              uint32
	easf06              uint32
	easf07              uint32
	easf08              uint32
	easf09              uint32
	easf10              uint32
	easd01              uint32
	easd02              uint32
	easd03              uint32
	easd04              uint32
	eait                uint32
	erq1                uint32
	erq2                uint32
	erq3                uint32
	erq4                uint32
	irms1               uint16
	irms2               uint16
	irms3               uint16
	urms1               uint16
	urms2               uint16
	urms3               uint16
	pref                uint8
	pcoup               uint8
	sinsts              uint16
	sinsts1             int16
	sinsts2             int16
	sinsts3             int16
	sinsti              uint16
	stge                string
	dpm1                string
	dpm1_timestamp      string
	fpm1                string
	fpm1_timestamp      string
	dpm2                string
	dpm2_timestamp      string
	fpm2                string
	fpm2_timestamp      string
	dpm3                string
	dpm3_timestamp      string
	fpm3                string
	fpm3_timestamp      string
	msg1                string
	msg2                string
	prm                 string
	relais              uint8
	ntarf               string
	njourf              string
	njourf_1            string
	pjourf_1            string
	ppointe             string
}

// NewLinkyCollector method to construct LinkyCollector
func NewLinkyCollector(device string, baudRate int, frameSize byte, parity serial.Parity, stopBits serial.StopBits) *LinkyCollector {
	return &LinkyCollector{
		device:    device,
		baudRate:  baudRate,
		frameSize: frameSize,
		parity:    parity,
		stopBits:  stopBits,
		linky_info: prometheus.NewDesc("linky_info",
			"Informations textuelles du compteur",
			[]string{"prm", "adsc", "vtic", "date", "ngtf", "ltarf", "stge", "msg1", "msg2", "ntarf"}, nil,
		),
		linky_index: prometheus.NewDesc("linky_index_watthours_total",
			"Index en Wh",
			[]string{"prm", "index"}, nil,
		),
		linky_current: prometheus.NewDesc("linky_current_amperes",
			"Courant efficace en A",
			[]string{"prm", "phase"}, nil,
		),
		linky_voltage: prometheus.NewDesc("linky_voltage_volts",
			"Tension efficace en V",
			[]string{"prm", "phase"}, nil,
		),
		linky_subscribed_power: prometheus.NewDesc("linky_subscribed_power_voltamperes",
			"Puissance apparente souscrite en VA",
			[]string{"prm", "type"}, nil,
		),
		linky_power: prometheus.NewDesc("linky_power_voltamperes",
			"Puissance apparente en VA",
			[]string{"prm", "direction", "phase"}, nil,
		),
		linky_load_management: prometheus.NewDesc("linky_load_management_info",
			"Informations relatives aux pointes mobiles",
			[]string{"prm", "dpm1", "dpm1_timestamp", "fpm1", "fpm1_timestamp", "dpm2", "dpm2_timestamp", "fpm2", "fpm2_timestamp", "dpm3", "dpm3_timestamp", "fpm3", "fpm3_timestamp", "pm_profile"}, nil,
		),
		linky_relays: prometheus.NewDesc("linky_relays",
			"État des relais, 0 = ouvert 1 = fermé, le premier est réel",
			[]string{"prm", "relay"}, nil,
		),
		linky_provider_day: prometheus.NewDesc("linky_provider_day_info",
			"Numéro du jour en cours, du prochain jour et de son profil",
			[]string{"prm", "current_day", "next_day", "next_day_profile"}, nil,
		),
	}
}

// Describe implements required describe function for all prometheus collectors
func (collector *LinkyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.linky_info
	ch <- collector.linky_index
	ch <- collector.linky_current
	ch <- collector.linky_voltage
	ch <- collector.linky_subscribed_power
	ch <- collector.linky_power
	ch <- collector.linky_load_management
	ch <- collector.linky_relays
	ch <- collector.linky_provider_day
}

// Collect implements required collect function for all prometheus collectors
func (collector *LinkyCollector) Collect(ch chan<- prometheus.Metric) {
	//for each descriptor or call other functions that do so.
	//Implement logic here to determine proper metric value to return to prometheus
	values := linkyValues{}
	err := collector.readSerial(&values)

	if err == nil {
		//Write latest value for each metric in the prometheus metric channel.
		//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
		ch <- prometheus.MustNewConstMetric(collector.linky_info, prometheus.GaugeValue, 1, values.prm, values.adsc, values.vtic, values.date, values.ngtf, values.ltarf, values.stge, values.msg1, values.msg2, values.ntarf)
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.east), values.prm, "east")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easf01), values.prm, "easf01")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easf02), values.prm, "easf02")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easf03), values.prm, "easf03")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easf04), values.prm, "easf04")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easf05), values.prm, "easf05")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easf06), values.prm, "easf06")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easf07), values.prm, "easf07")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easf08), values.prm, "easf08")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easf09), values.prm, "easf09")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easf10), values.prm, "easf10")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easd01), values.prm, "easd01")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easd02), values.prm, "easd02")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easd03), values.prm, "easd03")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.easd04), values.prm, "easd04")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.eait), values.prm, "eait")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.erq1), values.prm, "erq1")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.erq2), values.prm, "erq2")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.erq3), values.prm, "erq3")
		ch <- prometheus.MustNewConstMetric(collector.linky_index, prometheus.CounterValue, float64(values.erq4), values.prm, "erq4")
		ch <- prometheus.MustNewConstMetric(collector.linky_current, prometheus.GaugeValue, float64(values.irms1), values.prm, "1")
		ch <- prometheus.MustNewConstMetric(collector.linky_current, prometheus.GaugeValue, float64(values.irms2), values.prm, "2")
		ch <- prometheus.MustNewConstMetric(collector.linky_current, prometheus.GaugeValue, float64(values.irms3), values.prm, "3")
		ch <- prometheus.MustNewConstMetric(collector.linky_voltage, prometheus.GaugeValue, float64(values.urms1), values.prm, "1")
		ch <- prometheus.MustNewConstMetric(collector.linky_voltage, prometheus.GaugeValue, float64(values.urms2), values.prm, "2")
		ch <- prometheus.MustNewConstMetric(collector.linky_voltage, prometheus.GaugeValue, float64(values.urms3), values.prm, "3")
		ch <- prometheus.MustNewConstMetric(collector.linky_subscribed_power, prometheus.GaugeValue, float64(values.pref) * 1000, values.prm, "pref")
		ch <- prometheus.MustNewConstMetric(collector.linky_subscribed_power, prometheus.GaugeValue, float64(values.pcoup) * 1000, values.prm, "pcoup")
		ch <- prometheus.MustNewConstMetric(collector.linky_power, prometheus.GaugeValue, float64(values.sinsts), values.prm, "drawn", "sum")
		ch <- prometheus.MustNewConstMetric(collector.linky_power, prometheus.GaugeValue, float64(values.sinsts1), values.prm, "drawn", "1")
		ch <- prometheus.MustNewConstMetric(collector.linky_power, prometheus.GaugeValue, float64(values.sinsts2), values.prm, "drawn", "2")
		ch <- prometheus.MustNewConstMetric(collector.linky_power, prometheus.GaugeValue, float64(values.sinsts3), values.prm, "drawn", "3")
		ch <- prometheus.MustNewConstMetric(collector.linky_power, prometheus.GaugeValue, float64(values.sinsti), values.prm, "injected", "sum")
		ch <- prometheus.MustNewConstMetric(collector.linky_load_management, prometheus.GaugeValue, 1, values.prm, values.dpm1, values.dpm1_timestamp, values.fpm1, values.fpm1_timestamp, values.dpm2, values.dpm2_timestamp, values.fpm2, values.fpm2_timestamp, values.dpm3, values.dpm3_timestamp, values.fpm3, values.fpm3_timestamp, values.ppointe)
		ch <- prometheus.MustNewConstMetric(collector.linky_relays, prometheus.GaugeValue, float64(values.relais), values.prm, "relays")
		ch <- prometheus.MustNewConstMetric(collector.linky_provider_day, prometheus.GaugeValue, 1, values.prm, values.njourf, values.njourf_1, values.pjourf_1)
	} else {
		log.Errorf("Unable to read telemetry information : %s", err)
	}
}

// Read information from serial port
func (collector *LinkyCollector) readSerial(linkyValues *linkyValues) error {
	c := &serial.Config{Name: collector.device, Baud: collector.baudRate, Size: collector.frameSize, Parity: collector.parity, StopBits: collector.stopBits}
	stream, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(stream)
	started := false
	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			return err
		}

		line := string(bytes)

		// End loop when block ended
		if started && strings.Contains(line, string(0x03)) {
			break
		}

		// Start reading data when block started
		if strings.Contains(line, string(0x02)) {
			started = true
		}

		// Collect data
		if started {
			collector.proceedLine(linkyValues, line)
		}
	}
	return nil
}

// Proceed line by line information
func (collector *LinkyCollector) proceedLine(linkyValues *linkyValues, line string) {
	data := strings.Split(line, string(0x09))

	switch strings.ToLower(data[0]) {
	case "adsc":
		linkyValues.adsc = string(data[1])
	case "vtic":
		linkyValues.vtic = string(data[1])
	case "date":
		linkyValues.date = string(data[1])
	case "ngtf":
		linkyValues.ngtf = string(data[1])
	case "ltarf":
		linkyValues.ltarf = string(data[1])
	case "east":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.east = uint32(val)
	case "easf01":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easf01 = uint32(val)
	case "easf02":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easf02 = uint32(val)
	case "easf03":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easf03 = uint32(val)
	case "easf04":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easf04 = uint32(val)
	case "easf05":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easf05 = uint32(val)
	case "easf06":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easf06 = uint32(val)
	case "easf07":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easf07 = uint32(val)
	case "easf08":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easf08 = uint32(val)
	case "easf09":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easf09 = uint32(val)
	case "easf10":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easf10 = uint32(val)
	case "easd01":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easd01 = uint32(val)
	case "easd02":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easd02 = uint32(val)
	case "easd03":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easd03 = uint32(val)
	case "easd04":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.easd04 = uint32(val)
	case "eait":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.eait = uint32(val)
	case "erq1":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.erq1 = uint32(val)
	case "erq2":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.erq2 = uint32(val)
	case "erq3":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.erq3 = uint32(val)
	case "erq4":
		val, _ := strconv.ParseUint(data[1], 10, 32)
		linkyValues.erq4 = uint32(val)
	case "irms1":
		val, _ := strconv.ParseUint(data[1], 10, 16)
		linkyValues.irms1 = uint16(val)
	case "irms2":
		val, _ := strconv.ParseUint(data[1], 10, 16)
		linkyValues.irms2 = uint16(val)
	case "irms3":
		val, _ := strconv.ParseUint(data[1], 10, 16)
		linkyValues.irms3 = uint16(val)
	case "urms1":
		val, _ := strconv.ParseUint(data[1], 10, 16)
		linkyValues.urms1 = uint16(val)
	case "urms2":
		val, _ := strconv.ParseUint(data[1], 10, 16)
		linkyValues.urms2 = uint16(val)
	case "urms3":
		val, _ := strconv.ParseUint(data[1], 10, 16)
		linkyValues.urms3 = uint16(val)
	case "pref":
		val, _ := strconv.ParseUint(data[1], 10, 8)
		linkyValues.pref = uint8(val)
	case "pcoup":
		val, _ := strconv.ParseUint(data[1], 10, 8)
		linkyValues.pcoup = uint8(val)
	case "sinsts":
		val, _ := strconv.ParseUint(data[1], 10, 16)
		linkyValues.sinsts = uint16(val)
	case "sinsts1":
		val, _ := strconv.ParseInt(data[1], 10, 16)
		linkyValues.sinsts1 = int16(val)
	case "sinsts2":
		val, _ := strconv.ParseInt(data[1], 10, 16)
		linkyValues.sinsts2 = int16(val)
	case "sinsts3":
		val, _ := strconv.ParseInt(data[1], 10, 16)
		linkyValues.sinsts3 = int16(val)
	case "sinsti":
		val, _ := strconv.ParseUint(data[1], 10, 16)
		linkyValues.sinsti = uint16(val)
	case "stge":
		linkyValues.stge = string(data[1])
	case "dpm1":
		linkyValues.dpm1 = string(data[2])
		linkyValues.dpm1_timestamp = string(data[1])
	case "fpm1":
		linkyValues.fpm1 = string(data[2])
		linkyValues.fpm1_timestamp = string(data[1])
	case "dpm2":
		linkyValues.dpm2 = string(data[2])
		linkyValues.dpm2_timestamp = string(data[1])
	case "fpm2":
		linkyValues.fpm2 = string(data[2])
		linkyValues.fpm2_timestamp = string(data[1])
	case "dpm3":
		linkyValues.dpm3 = string(data[2])
		linkyValues.dpm3_timestamp = string(data[1])
	case "fpm3":
		linkyValues.fpm3 = string(data[2])
		linkyValues.fpm3_timestamp = string(data[1])
	case "msg1":
		linkyValues.msg1 = string(data[1])
	case "msg2":
		linkyValues.msg2 = string(data[1])
	case "prm":
		linkyValues.prm = string(data[1])
	case "relais":
		val, _ := strconv.ParseUint(data[1], 10, 8)
		linkyValues.relais = uint8(val)
	case "ntarf":
		linkyValues.ntarf = string(data[1])
	case "njourf":
		linkyValues.njourf = string(data[1])
	case "njourf+1":
		linkyValues.njourf_1 = string(data[1])
	case "pjourf+1":
		linkyValues.pjourf_1 = string(data[1])
	case "ppointe":
		linkyValues.ppointe = string(data[1])
	}
}
