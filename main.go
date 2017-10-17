package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/lovoo/ipmi_exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

var (
	listenAddress = flag.String("web.listen", ":9289", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.path", "/metrics", "Path under which to expose metrics.")
	ipmiBinary    = flag.String("ipmi.path", "ipmitool", "Path to ipmi binar")
	ipmiHost      = flag.String("ipmi.host", "", "IPMI server address")
	ipmiUser      = flag.String("ipmi.user", "", "IPMI server username")
	ipmiPass      = flag.String("ipmi.password", "", "IPMI server password")
	showVersion   = flag.Bool("version", false, "Show version information and exit")
)

func init() {
	prometheus.MustRegister(version.NewCollector("ipmi_exporter"))
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("ipmi_exporter"))
		os.Exit(0)
	}
 
	if *ipmiHost != "" {
		*ipmiBinary += " -H " + string(*ipmiHost)
	}

	if *ipmiUser != "" {
		*ipmiBinary += " -U" + string(*ipmiUser)
	}

	if *ipmiPass != "" {
		*ipmiBinary += " -P" + string(*ipmiPass)
	}

	//log.Infoln("ipmi command", *ipmiBinary)

	log.Infoln("Starting IPMI Exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	prometheus.MustRegister(collector.NewExporter(*ipmiBinary))

	handler := promhttp.Handler()
	if *metricsPath == "" || *metricsPath == "/" {
		http.Handle(*metricsPath, handler)
	} else {
		http.Handle(*metricsPath, handler)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>IPMI Exporter</title></head>
			<body>
			<h1>IPMI Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		})
	}

	log.Infoln("Listening on", *listenAddress, *metricsPath)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
