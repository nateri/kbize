package netgo

import (
	//"flag"
	"fmt"
	"os"
	//"regexp"
	//"time"
	//"strings"
	//"net/url"
	"net/http"
	//"encoding/text"
	"encoding/json"
	//"strconv"
	//"github.com/op/go-logging"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
)

//var log = logging.MustGetLogger("netgo")

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
//  Handlers

func RootRequestHandler(w *bytes.Buffer, req *http.Request) error {
	LogRequest("Root", req)

	if req.URL.Path != "/" {
		log.Error("Path not found!")
		return errors.New("NotFound")
	}

	fmt.Fprintf(w, "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n")
	fmt.Fprintf(w, "<html xmlns=\"http://www.w3.org/1999/xhtml\">\n")

	//var test_str string
	//criteria, _ := Netgo_GenerateHtml(req.URL)

	criteria2 := struct {
		Name   string
		Method string
		Host   string
	}{
		"NetgoTest!",
		"Magick",
		"",
	}

	if err := WriteHtmlFromTemplate(w, "html/head.html", criteria2); err != nil {
		return err
	}

	fmt.Fprintf(w, "<body>\n")

	fmt.Fprintf(w, "<p>Welcome<p>\n")

	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")

	return nil
}

type JsonUser struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Appointee bool   `json:"appointee"`
	Score     uint64 `json:"score"`
}
type KbizeUser struct {
	UUID string
	Name string
	KB   bool
}
type is_user_kb func(JsonUser) bool

func KbizeRequestHandler(w *bytes.Buffer, req *http.Request) error {
	LogRequest("Kbize", req)

	// @TODO: Get json data from DB
	path, err := filepath.Abs("/test/kbize_db.json")
	if err != nil {
		return err
	}

	in, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n")
	fmt.Fprintf(w, "<html xmlns=\"http://www.w3.org/1999/xhtml\">\n")

	//var test_str string
	//criteria, _ := Netgo_GenerateHtml(req.URL)

	criteria2 := struct {
		Name   string
		Method string
		Host   string
	}{
		"NetgoTest!",
		"Magick",
		"",
	}

	if err := WriteHtmlFromTemplate(w, "html/head.html", criteria2); err != nil {
		return err
	}

	fmt.Fprintf(w, "<body>\n")

	//x := &tar{foo: {bar: {value: 2}}}

	var users []JsonUser

	if err := json.Unmarshal(in, &users); err != nil {
		return err
	}

	if err := WriteHtmlFromTemplate(w, "html/users.html", users); err != nil {
		return err
	}

	var eddied_score uint64
	eddied_score = 50
	fmt.Fprintf(w, "<p><p><p>Generating Kbized Users for alg [Eddied] (Appointee or Score >= %d)<p>\n", eddied_score)

	eddied := func(user JsonUser) bool {
		if user.Appointee {
			return true
		}
		if user.Score >= eddied_score {
			return true
		}
		return false
	}
	users_kbize := GetKbized(users, eddied)

	if err := WriteHtmlFromTemplate(w, "html/eddied.html", users_kbize); err != nil {
		return err
	}

	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")

	return nil
}

func GetKbized(users []JsonUser, eval is_user_kb) []KbizeUser {
	kbize_users := make([]KbizeUser, len(users))

	for idx, user := range users {
		kbize_users[idx].UUID = user.UUID
		kbize_users[idx].Name = user.Name
		kbize_users[idx].KB = eval(user)
	}

	return kbize_users
}

func NetgetRequestHandler(w *bytes.Buffer, req *http.Request) error {
	LogRequest("Netget", req)

	fmt.Fprintf(w, "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n")
	fmt.Fprintf(w, "<html xmlns=\"http://www.w3.org/1999/xhtml\">\n")

	//var test_str string
	// define criteria
	criteria2 := struct {
		Test string
	}{
		"testing...",
	}
	log.Debug("%+v", criteria2)

	criteria, _ := Netgo_GenerateHtml(req.URL)

	if err := WriteHtmlFromTemplate(w, "html/head.html", criteria); err != nil {
		return err
	}
	fmt.Fprintf(w, "<body>\n\n")
	log.Info("%+v", criteria)

	if err := WriteHtmlFromTemplate(w, "html/query.html", criteria); err != nil {
		return err
	}

	//results := TestTpbYqlJson(TestTpbJsonStr)

	fmt.Fprintf(w, "<div id=\"SearchResults\">\n<div id=\"content\">\n<div id=\"main-content\">\n\n<table id=\"searchResult\">\n\n")

	if err := WriteHtmlFromTemplate(w, "html/resultsheader.html", criteria); err != nil {
		return err
	}
	//WriteHtmlFromTemplate(w, "html/resultsentry.html", results)

	fmt.Fprintf(w, "</table></div>\n\n")
	fmt.Fprintf(w, "</div></div></div><!-- //div:content -->\n\n")

	fmt.Fprintf(w, "</body>\n</html>\n\n")

	return nil
}

func TpbSearchRequestHandler(w *bytes.Buffer, req *http.Request) error {
	LogRequest("TpbSearch", req)

	req.Close = true
	devmode := false

	tpb_criteria, test_str := GetTpbSearchCriteriaFromUrl(req.URL)
	if test_str != "" {
		devmode = true
	}
	if devmode {
		Serialize("criteria.log", tpb_criteria)
	}

	results := SearchTpb(tpb_criteria)

	out := struct {
		Criteria TpbSearchCriteria
		Results  []SearchResult
	}{
		tpb_criteria,
		results,
	}

	response_json, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		Serialize("results.log", err)
		return err
	}

	resultDbg := struct {
		RemoteAddr string
		NumResults int
		//Time string
	}{
		req.RemoteAddr,
		len(results),
	}
	log.Info("[ResultDbg ] %v", resultDbg)

	if devmode {
		Serialize("results.log", string(response_json))
	}
	fmt.Fprintf(w, "%+v", string(response_json))

	return nil
}

func ShutdownRequestHandler(w *bytes.Buffer, req *http.Request) error {
	LogRequest("Shutdown", req)

	req.Close = true

	log.Critical("[Quit Gracefully]")
	os.Exit(1)

	return nil
}

func LogRequest(handler string, req *http.Request) {

	requestDbg := struct {
		Handler    string
		RequestURI string
		RemoteAddr string
		//Time
	}{
		handler,
		req.RequestURI,
		req.RemoteAddr,
		//time.Now().Local().String(),
	}

	log.Info("%v", requestDbg)
	SerializeAppend("request.log", requestDbg)
}

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
// @TODO Netget
func SaveFileFromUri(out_path string, in_path string) {
	out, err := os.Create(out_path)
	defer out.Close()
	if err != nil {
		log.Info("[File not created: %+v]", err)
		return
	}

	resp, err := http.Get("in_path")
	defer resp.Body.Close()
	if err != nil {
		log.Info("[Path not found: %+v]", err)
		return
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Info("[Copy failed: %+v]", err)
		return
	}
}
