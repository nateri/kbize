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
	"errors"
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

	http.Handle("/shutdown/", DynHandle{netgo.ShutdownRequestHandler, "Shutdown", false})

	// Mandatory root-based resources
	http.HandleFunc("/sitemap.xml", FileRequestHandler)
	http.HandleFunc("/favicon.ico", FileRequestHandler)
	http.HandleFunc("/robots.txt", FileRequestHandler)

}

type Headers struct {
	Origin    string
	TodoLater bool
}
type ErrorPage struct {
	content string
}

func (p *ErrorPage) Set(in string) {
	p.content = in
}

//func (b *Buffer) Bytes() []byte
func (p *ErrorPage) Bytes() []byte {
	var page bytes.Buffer

	if err := WriteErrorPageStart(&page); err != nil {
		// @TODO: What do we do when the error page itself fails to be generated?!
		// Maybe this should be hard-coded rather than read from a file...
		return nil
	}

	page.Write([]byte(p.content))

	WriteErrorPageEnd(&page)

	return page.Bytes()
}

func (d DynHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("[ServeHTTP] [%s] [Add Headers: %t]", d.dbg, d.headers)

	// Handle error gracefully
	recover_handler := func() {
		// w went into a bad state during write
		//  Reset and return
		//w.Reset()
	}
	// Unused -- ServeHTTP cannot return an error
	var unused_err error
	defer Recoverable(&unused_err, recover_handler)

	var page bytes.Buffer
	var error_page ErrorPage

	var headers Headers
	headers.TodoLater = d.headers
	headers.Origin = r.Header.Get("Origin")
	set_headers(w, headers)

	// GracefulHandler should NOT write to w unless there is no error
	if err := d.gh(&page, r); err != nil {
		log.Debug("[err] [%s]", d.dbg)

		switch err.Error() {
		case "NotFound":
			log.Warning("[ServeHTTP] [NotFound]")

			w.WriteHeader(http.StatusNotFound)
			error_page.Set(val("404_error"))

		default:
			log.Warning("[ServeHTTP] [start]\n" + page.String())
			log.Warning("[ServeHTTP] [end]")

			w.WriteHeader(http.StatusInternalServerError)
			error_page.Set(err.Error())
		}

		w.Write(error_page.Bytes())
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

func set_headers(w http.ResponseWriter, headers Headers) {

	log.Debug("[set_headers]")

	if headers.Origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", headers.Origin)
	}

	if headers.TodoLater {
		//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}
func WriteErrorPageStart(w *bytes.Buffer) error {

	// does it make sense to catch errors here?
	if err := recoverable_write(w, val("html_tag_start")); err != nil {
		return err
	}

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

	// does it make sense to catch errors here?
	if err := recoverable_write(w, val("body_tag_start")); err != nil {
		return err
	}
	if err := recoverable_write(w, val("generic_error")); err != nil {
		return err
	}

	return nil
}
func WriteErrorPageEnd(w *bytes.Buffer) {
	fmt.Fprintf(w, val("body_tag_end"))
	fmt.Fprintf(w, val("html_tag_end"))
}

type ResetHandler func()

func Recoverable(err *error, h ResetHandler) {

	if e := recover(); e != nil {
		// e is the interface{} typed-value we passed to panic()
		fmt.Println("Whoops: ", e) // Prints "Whoops: <err>!"

		// find out exactly what the error was and set err
		switch x := e.(type) {
		case string:
			*err = errors.New(x)
		case error:
			*err = x
		default:
			*err = errors.New("Unknown panic")
		}
		// return the modified err

		h()
	}
}
func recoverable_write(w *bytes.Buffer, instr string) (err error) {

	fmt.Fprintf(w, instr)
	// Fprintf may panic --
	//  recover() will return immediately on fail
	return nil
}

func val(Name string) string {

	switch Name {

	case "html_tag_start":
		return "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n<html xmlns=\"http://www.w3.org/1999/xhtml\">\n"
	case "html_tag_end":
		return "</html>\n"

	case "body_tag_start":
		return "<body>\n"
	case "body_tag_end":
		return "</body>\n"

	case "404_error":
		return "404 page not found"

	default:
	}

	return ""
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

/*
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Catastrphic error generating the error page!")
	fmt.Fprintln(w, err.Error())
*/
