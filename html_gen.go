package netgo

import (
	//"bytes"
	//"flag"
	//"fmt"
	//stdlog "log"
	"html/template"
	//"io"
	//"io/ioutil"
	//"os"
	//"regexp"
	//"time"
	//"net/http"
	"bytes"
	//"net/url"
	//"strings"
	//"encoding/xml"
	//"encoding/json"
	//"strconv"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("netgo")

var (
// a source of numbers, for naming temporary files
//uniq = make(chan int)
)

var (
//portMapRule = regexp.MustCompile(`^([1-9][0-9]*?):([1-9][0-9]*?)$`)
//timeout = flag.Uint("timeout", 10, "timeout sec to wait between search queries")
)

type MethodType string

type IMethodType interface {
	Named() string
}

type NetwebHtmlCriteria struct {
	Name string

	Method MethodType

	//Order []OrderType // Can only be len<=6 and should only have 1 of each type
	//Order OrderList
	Orders map[string]string

	// Request 1 page at a time...
	Page string

	CategoryChecks map[string]string
}

var CategoryAbv = [...]string{
	"0",
	"A",
	"V",
	"P",
	"G",
	"N",
	"O",
}
var CategoryStrings = [...]string{
	"None",
	"Audio",
	"Video",
	"Apps",
	"Games",
	"Nsfw",
	"Other",
}

var OrderAbv = [...]string{
	"0",
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
	"A",
	"B",
	"C",
}

var OrderMap = map[string]string{
	"0": "Unordered",
	"1": "Name A-Z",
	"2": "Name Z-A",
	"3": "Date New",
	"4": "Date Old",
	"5": "Size Big",
	"6": "Size Small",
	"7": "Seed Most",
	"8": "Seed Least",
	"9": "Leech Most",
	"A": "Leech Least",
	"B": "Category A-Z",
	"C": "Category Z-A",
}

func get_orders_from_abvs(abv []string) []string {
	var orders []string
	for _, o := range abv {
		orders = append(orders, OrderMap[o])
	}
	return orders
}

type NetgoHtml struct {
	Host string
	This string
	Path string

	Search     string
	Method     string
	Page       string
	Categories struct {
		Short   string
		Checked map[string]string
	}
	Orders struct {
		Short         string
		SelectedAbv   []string
		UnselectedAbv []string
		Selected      []string
		Unselected    []string
	}
}

func init_category_checks(s string) map[string]string {
	checks := make(map[string]string)
	for _, cat := range CategoryAbv {
		checks[cat] = s
	}
	/*
		for _, checks := range CategoryStrings {
			checks[category] = s
		}
	*/
	return checks
}

func update_category_checks(s string) map[string]string {

	if s == "" {
		return init_category_checks("checked")
	}

	checks := init_category_checks("")

	for _, c := range s {
		checks[string(c)] = "checked"
	}

	return checks
}

func is_sort_abv(abv string) bool {
	// "None" is not a choice
	for _, e := range OrderAbv[1:] {
		if abv == e {
			return true
		}
	}
	return false
}
func update_selected_sorts(s string) []string {
	var selected []string

	for _, c := range s {
		abv := string(c)
		if is_sort_abv(abv) {
			selected = append(selected, abv)
		}
	}

	return selected
}

func update_available_sorts(selected []string) []string {
	var available []string

	var skip_name, skip_date, skip_size, skip_seed, skip_leech, skip_category bool
	for _, entry := range selected {
		switch entry {
		case OrderAbv[1], OrderAbv[2]:
			skip_name = true
		case OrderAbv[3], OrderAbv[4]:
			skip_date = true
		case OrderAbv[5], OrderAbv[6]:
			skip_size = true
		case OrderAbv[7], OrderAbv[8]:
			skip_seed = true
		case OrderAbv[9], OrderAbv[10]:
			skip_leech = true
		case OrderAbv[11], OrderAbv[12]:
			skip_category = true
		}
	}
	if !skip_name {
		available = append(available, OrderAbv[1:3]...)
	}
	if !skip_date {
		available = append(available, OrderAbv[3:5]...)
	}
	if !skip_size {
		available = append(available, OrderAbv[5:7]...)
	}
	if !skip_seed {
		available = append(available, OrderAbv[7:9]...)
	}
	if !skip_leech {
		available = append(available, OrderAbv[9:11]...)
	}
	if !skip_category {
		available = append(available, OrderAbv[11:13]...)
	}

	return available
}

func WriteHtmlFromTemplate(w *bytes.Buffer, in_template string, data interface{}) error {
	t, err := template.ParseFiles(in_template)
	if err != nil {
		log.Debug("[Parse failed: %+v]", err)
		return err
	}

	Serialize("test.log", t)

	err = t.Execute(w, data)
	if err != nil {
		log.Debug("[Execute failed: %+v]", err)
		return err
	}

	return nil
}
