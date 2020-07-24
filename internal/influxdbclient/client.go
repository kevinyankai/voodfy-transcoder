package influxdbclient

import (
	"fmt"
	"time"

	"github.com/Voodfy/voodfy-transcoder/internal/settings"
	influxdb2 "github.com/influxdata/influxdb-client-go"
)

// Influx struct used to mantain the influx client
type Influx struct {
	Client influxdb2.Client
}

// NewClient return a instance of InfluxDB
func NewClient() Influx {
	client := influxdb2.NewClientWithOptions(fmt.Sprintf("http://%s", settings.InfluxdbSetting.Host), settings.InfluxdbSetting.Password,
		influxdb2.DefaultOptions().SetBatchSize(20))
	defer client.Close()
	return Influx{Client: client}
}

// Send metrics to influx
func (i *Influx) Send(t interface{}, field, value string) {
	writeAPI := i.Client.WriteApi(settings.InfluxdbSetting.User, settings.InfluxdbSetting.DB)
	p := influxdb2.NewPointWithMeasurement("transcoder").
		AddTag("stat:task", fmt.Sprintf("resource_id_%v", t)).
		AddField(field, value).
		SetTime(time.Now())
	writeAPI.WritePoint(p)
	writeAPI.Flush()
}
