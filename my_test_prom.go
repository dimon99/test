package main

import (
	"net/http"
	"os/exec"
	"github.com/prometheus/client_golang/prometheus"
	"fmt"
	"strconv"
	"regexp"
	"github.com/gorilla/mux"
	"time"
	"log"
	"github.com/jasonlvhit/gocron"
	"os"
	"net"
)

func prometeushendler() http.Handler {
	//_,s := http.Get("http://localhost:8080/info")
	//fmt.Println("hhh",s)
	return prometheus.Handler()

}

func mydata(options prometheus.Gauge) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		out, _ := exec.Command("/bin/sh", "-c", "ps | wc -l | awk '{ print $1}'").Output()
		re := regexp.MustCompile("[0-9]+")
		c := re.FindAllString(string(out),-1)
		fmt.Println(c[0])
		mm,_ := strconv.ParseFloat(c[0], 64)
		fmt.Println(mm)
		options.Set(mm)
		w.Write(out)
	}
}

func hellow(w http.ResponseWriter, req *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "Overview\n")
        fmt.Fprintln(w, "                     ")
        fmt.Fprintln(w, "Hellow world! This is Prometeus plugin for TEST1 purposes!")
}

func task() {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			//fmt.Println(ipv4)
			ip := fmt.Sprintf("%s",ipv4)
			_,_ = http.Get("http://"+ip+":8181/info", )
			//fmt.Println(r,e)
			
		}
	}
	//_,_ = http.Get("http://localhost:8181/info")
	//fmt.Println(r)
}
func main() {

	cpuTemp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "proc_num_count",
		Help: "Current file info",
	})
	r := mux.NewRouter()
	prometheus.Register(cpuTemp)
	s :=&http.Server{
		Addr: ":8181",
		ReadTimeout: 8 * time.Second,
		WriteTimeout: 8 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler: r,
	}
	gocron.Every(1).Second().Do(task)
	gocron.Start()

	r.Handle("/metrics", prometeushendler())
	r.Handle("/info", mydata(cpuTemp)).Methods("GET")
        r.HandleFunc("/", hellow).Methods("GET")
	log.Fatal(s.ListenAndServe())
}
