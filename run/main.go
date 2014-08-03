package main

import (
	//"bytes"
	"flag"
	//"fmt"
	//"log/syslog"
	stdlog "log"
	//"html/template"
	//"io"
	//"io/ioutil"
	"os"
	//"regexp"
	//"time"
	//"strings"
	"net/http"
	//"net/url"
	//"encoding/xml"
	//"encoding/json"
	//"strconv"
	"bytes"
	"fmt"
	"github.com/nateri/kbize"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("netgo")

var (
// a source of numbers, for naming temporary files
//uniq = make(chan int)
)

var (
	httpListen = flag.String("http", ":3999", "host:port to listen on")
	//htmlOutput = flag.Bool("html", false, "render program output as HTML")

	//portMapRule = regexp.MustCompile(`^([1-9][0-9]*?):([1-9][0-9]*?)$`)
	//timeout = flag.Uint("timeout", 10, "timeout sec to wait between search queries")
)

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
// Main

func main() {
	log.Critical("[main entered]")

	InitLogging()

	flag.Parse()

	//netgo.TestNetget()

	//timeout := time.Duration(*timeout) * time.Second

	AddHttpHandlers()

	log.Critical("[Starting Service] [%s]", *httpListen)
	log.Fatal(http.ListenAndServe(*httpListen, nil))
}

type GracefulHandler func(*bytes.Buffer, *http.Request) error

type DynHandle struct {
	gh      GracefulHandler
	dbg     string
	headers bool
}

func AddHttpHandlers() {
	log.Debug("[Adding Http RequestHandlers]")

	http.Handle("/", DynHandle{netgo.RootRequestHandler, "Root", true})

	http.HandleFunc("/static/", StaticRequestHandler)

	http.Handle("/kbize/", DynHandle{netgo.KbizeRequestHandler, "Kbize", false})

	http.Handle("/netget/", DynHandle{netgo.KbizeRequestHandler, "Netget", false})
	//http.Handle("/tpb/", DynHandle{netgo.TpbSearchRequestHandler, "TpbSearch", true})

	http.Handle("/shutdown/", DynHandle{netgo.ShutdownRequestHandler, "Shutdown", false})

	// Mandatory root-based resources
	http.HandleFunc("/sitemap.xml", FileRequestHandler)
	http.HandleFunc("/favicon.ico", FileRequestHandler)
	http.HandleFunc("/robots.txt", FileRequestHandler)

}

func (d DynHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("[ServeHTTP] [%s] [Add Headers: %t]", d.dbg, d.headers)

	var page bytes.Buffer
	var error_page bytes.Buffer

	if d.headers {
		log.Debug("[ServeHTTP] [write headers]")

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if err := WriteErrorPageStart(&error_page); err != nil {
		// @TODO: What do we do when the error page itself fails to be generated?!
		// Maybe this should be hard-coded rather than read from a file...

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Catastrphic error generating the error page!")
		fmt.Fprintln(w, err.Error())
		return
	}

	// GracefulHandler should NOT write to w unless there is no error
	if err := d.gh(&page, r); err != nil {
		log.Debug("[err] [%s]", d.dbg)

		switch err.Error() {
		case "NotFound":
			log.Warning("[ServeHTTP] [NotFound]")

			//w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(&error_page, "404 page not found")

			WriteErrorPageEnd(&error_page)
			fmt.Fprintln(w, &error_page)

		default:
			log.Warning("[ServeHTTP] [start]\n" + page.String())
			log.Warning("[ServeHTTP] [end]")

			//w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(&error_page, err.Error())

			WriteErrorPageEnd(&error_page)
			fmt.Fprintln(w, &error_page)
		}

		return
	}

	w.Write(page.Bytes())
}

func StaticRequestHandler(w http.ResponseWriter, req *http.Request) {
	netgo.LogRequest("Static", req)
	// @TODO: ServeFile writes directly to ResponseWriter
	http.ServeFile(w, req, req.URL.Path[1:])
}
func FileRequestHandler(w http.ResponseWriter, req *http.Request) {
	netgo.LogRequest("File", req)
	// @TODO: ServeFile writes directly to ResponseWriter
	http.ServeFile(w, req, "./static/img"+req.URL.Path)
}

func InitLogging() {
	// Customize the output format
	logging.SetFormatter(logging.MustStringFormatter("▶ %{level:.1s} %{message}"))

	// Setup one stdout and one syslog backend.
	console_log := logging.NewLogBackend(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
	console_log.Color = false

	//var file_log_filter = syslog.LOG_DEBUG|syslog.LOG_LOCAL0|syslog.LOG_CRIT
	//log.Info("[%d]", file_log_filter)
	//file_log := logging.NewSyslogBackend("")
	//file_log.Color = false

	// Combine them both into one logging backend.
	logging.SetBackend(console_log)

	logging.SetLevel(logging.INFO, "netgo")
}

func WriteErrorPageStart(w *bytes.Buffer) error {

	fmt.Fprintf(w, "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n")
	fmt.Fprintf(w, "<html xmlns=\"http://www.w3.org/1999/xhtml\">\n")

	head := struct {
		Name   string
		Method string
		Host   string
	}{
		"NetgoTest!",
		"Errors",
		"",
	}

	if err := netgo.WriteHtmlFromTemplate(w, "html/head.html", head); err != nil {
		return err
	}

	fmt.Fprintf(w, "<body>\n")
	return nil
}
func WriteErrorPageEnd(w *bytes.Buffer) {
	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")
}
