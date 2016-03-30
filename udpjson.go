// +build !windows

package udpjson

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))
	return []byte(stamp), nil
}

type JSONCounter struct {
	Time   JSONTime `json:"time"`
	Type   string   `json:"type"`
	Metric string   `json:"metric"`
	Value  int64    `json:"value"`
}

type JSONGauge struct {
	Time   JSONTime `json:"time"`
	Type   string   `json:"type"`
	Metric string   `json:"metric"`
	Value  int64    `json:"value"`
}

type JSONGaugeFloat64 struct {
	Time   JSONTime `json:"time"`
	Type   string   `json:"type"`
	Metric string   `json:"metric"`
	Value  float64  `json:"value"`
}

type JSONHealthcheck struct {
	Time   JSONTime `json:"time"`
	Type   string   `json:"type"`
	Metric string   `json:"metric"`
	Error  string   `json:"string"`
}

type JSONHistogram struct {
	Time    JSONTime `json:"time"`
	Type    string   `json:"type"`
	Metric  string   `json:"metric"`
	Count   int64    `json:"count"`
	Min     int64    `json:"minimum"`
	Max     int64    `json:"maximum"`
	Mean    float64  `json:"mean"`
	StdDev  float64  `json:"stddev"`
	Pct50   float64  `json:"pct50"`
	Pct75   float64  `json:"pct75"`
	Pct95   float64  `json:"pct95"`
	Pct99   float64  `json:"pct99"`
	Pct9999 float64  `json:"pct9999"`
}

type JSONMeter struct {
	Time     JSONTime `json:"time"`
	Type     string   `json:"type"`
	Metric   string   `json:"metric"`
	Count    int64    `json:"count"`
	Rate1m   float64  `json:"rate1m"`
	Rate5m   float64  `json:"rate5m"`
	Rate15m  float64  `json:"rate15m"`
	RateMean float64  `json:"ratemean"`
}

type JSONTimer struct {
	Time     JSONTime `json:"time"`
	Type     string   `json:"type"`
	Metric   string   `json:"metric"`
	Count    int64    `json:"count"`
	Min      int64    `json:"minimum"`
	Max      int64    `json:"maximum"`
	Mean     float64  `json:"mean"`
	StdDev   float64  `json:"stddev"`
	Pct50    float64  `json:"pct50"`
	Pct75    float64  `json:"pct75"`
	Pct95    float64  `json:"pct95"`
	Pct99    float64  `json:"pct99"`
	Pct9999  float64  `json:"pct9999"`
	Rate1m   float64  `json:"rate1m"`
	Rate5m   float64  `json:"rate5m"`
	Rate15m  float64  `json:"rate15m"`
	RateMean float64  `json:"ratemean"`
}

// Output each metric in the given registry to syslog periodically using
// the given syslogger.
func UDPJSON(r Registry, d time.Duration, s net.Conn) {
	for _ = range time.Tick(d) {
		now := JSONTime(time.Now())
		r.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case Counter:
				res1D := &JSONCounter{
					Time:   now,
					Type:   "counter",
					Metric: name,
					Value:  metric.Count(),
				}
				res1B, _ := json.Marshal(res1D)
				s.Write(res1B)
			case Gauge:
				res1D := &JSONGauge{
					Time:   now,
					Type:   "gauge",
					Metric: name,
					Value:  metric.Value(),
				}
				res1B, _ := json.Marshal(res1D)
				s.Write(res1B)
			case GaugeFloat64:
				res1D := &JSONGaugeFloat64{
					Time:   now,
					Type:   "gauge",
					Metric: name,
					Value:  metric.Value(),
				}
				res1B, _ := json.Marshal(res1D)
				s.Write(res1B)
			case Healthcheck:
				metric.Check()
				res1D := &JSONHealthcheck{
					Time:   now,
					Type:   "gauge",
					Metric: name,
					Error:  fmt.Sprintf("%v", metric.Error()),
				}
				res1B, _ := json.Marshal(res1D)
				s.Write(res1B)
			case Histogram:
				h := metric.Snapshot()
				ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				res1D := &JSONHistogram{
					Time:    now,
					Type:    "histogram",
					Metric:  name,
					Count:   h.Count(),
					Min:     h.Min(),
					Max:     h.Max(),
					Mean:    h.Mean(),
					StdDev:  h.StdDev(),
					Pct50:   ps[0],
					Pct75:   ps[1],
					Pct95:   ps[2],
					Pct99:   ps[3],
					Pct9999: ps[4],
				}
				res1B, _ := json.Marshal(res1D)
				s.Write(res1B)
			case Meter:
				m := metric.Snapshot()
				res1D := &JSONMeter{
					Time:     now,
					Type:     "meter",
					Metric:   name,
					Count:    m.Count(),
					Rate1m:   m.Rate1(),
					Rate5m:   m.Rate5(),
					Rate15m:  m.Rate15(),
					RateMean: m.RateMean(),
				}
				res1B, _ := json.Marshal(res1D)
				s.Write(res1B)
			case Timer:
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				res1D := &JSONTimer{
					Time:     now,
					Type:     "histogram",
					Metric:   name,
					Count:    t.Count(),
					Min:      t.Min(),
					Max:      t.Max(),
					Mean:     t.Mean(),
					StdDev:   t.StdDev(),
					Pct50:    ps[0],
					Pct75:    ps[1],
					Pct95:    ps[2],
					Pct99:    ps[3],
					Pct9999:  ps[4],
					Rate1m:   t.Rate1(),
					Rate5m:   t.Rate5(),
					Rate15m:  t.Rate15(),
					RateMean: t.RateMean(),
				}
				res1B, _ := json.Marshal(res1D)
				s.Write(res1B)
			}
		})
	}
}
