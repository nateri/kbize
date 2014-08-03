package netgo

import (
	//"strconv"
	"flag"
	"strings"
	//"log"
	"github.com/op/go-logging"
	"io/ioutil"
	"net/http"
	"net/url"
)

var log = logging.MustGetLogger("netgo")

var (
	TpbMethod = "yql_json"
)

type TpbSearchCriteria struct {
	Name     string
	Method   string
	Category string
	Order    string
	Page     string
}

func SearchTpbRaw(tpb_criteria TpbSearchCriteria) string {
	yql_format := "&format=json"
	return TpbSearchRaw(tpb_criteria, yql_format)
}

func SearchTpb(tpb_criteria TpbSearchCriteria) []SearchResult {
	yql_format := "&format=json"
	return TpbSearch(tpb_criteria, yql_format)
}

type TpbCategoryMapType map[CategoryType]string

var TpbCategoryMap = TpbCategoryMapType{
	Category_Uncategorized: "0",
	Category_Audio:         "100",
	Category_Audio_Other:   "199",
	Category_Video:         "200",
	Category_Video_Other:   "299",
	Category_Apps:          "300",
	//"App_Other": "399",
	Category_Games: "400",
	//"Games_Other": "499",
	Category_Nsfw: "500",
	//"Po_Other": "599",
	Category_Other: "600",
	//"Other_Other": "699",
}

type TpbOrderMapType map[OrderType]string

var TpbOrderMap = TpbOrderMapType{
	Order_None:        "0",
	Order_Name_A:      "1",
	Order_Name_Z:      "2",
	Order_Date_New:    "3",
	Order_Date_Old:    "4",
	Order_Size_Big:    "5",
	Order_Size_Small:  "6",
	Order_Seed_Most:   "7",
	Order_Seed_Least:  "8",
	Order_Leech_Most:  "9",
	Order_Leech_Least: "10",
	Order_Category_A:  "13",
	Order_Category_Z:  "14",
	Order_Unordered:   "99",
}

type TpbMethodMapType map[MethodType]string

var TpbMethodMap = TpbMethodMapType{
	Method_TpbFile: "search",
	Method_TpbUser: "user",
}

func ConvertCriteriaForTpbCategory(in_category CategoryType) string {
	return TpbCategoryMap[in_category]
}
func ConvertCriteriaForTpbOrder(in_order OrderType) string {
	return TpbOrderMap[in_order]
}
func ConvertCriteriaForTpbMethod(in_method MethodType) string {
	return TpbMethodMap[in_method]
}

func ConvertCriteriaForTpb(in_criteria SearchCriteria) TpbSearchCriteria {
	var tpb_criteria TpbSearchCriteria

	tpb_criteria.Name = in_criteria.Name
	tpb_criteria.Method = ConvertCriteriaForTpbMethod(in_criteria.Method)
	tpb_criteria.Category = ConvertCriteriaForTpbCategory(in_criteria.Category)
	tpb_criteria.Order = ConvertCriteriaForTpbOrder(in_criteria.Order)
	tpb_criteria.Page = in_criteria.Page

	return tpb_criteria
}

func GetHttpResponseAsString(http_request string) string {
	res, err := http.Get(http_request)
	if err != nil {
		log.Fatal(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(result)
}

func TpbGetResults(yql_response string, yql_format string) []SearchResult {
	if yql_format != "&format=json" {
		return nil
	}

	return TpbGetResultsJson(yql_response)
}

func TpbSearchRaw(criteria TpbSearchCriteria, yql_format string) string {

	if criteria.Name == "" {
		log.Info("[TpbSearchRaw] [abort: empty search]")
		return ""
	}

	yql_url := TpbCreateRequestUrlUsingYql(criteria, yql_format)
	log.Debug("[yql_url]\n%s", yql_url)

	yql_response := GetHttpResponseAsString(yql_url)

	return yql_response
}

func TpbSearch(criteria TpbSearchCriteria, yql_format string) []SearchResult {

	if criteria.Name == "" {
		log.Info("[TpbSearch] [abort: empty search]")
		return []SearchResult{}
	}

	yql_url := TpbCreateRequestUrlUsingYql(criteria, yql_format)
	log.Debug("[yql_url]\n%s", yql_url)

	yql_response := GetHttpResponseAsString(yql_url)
	log.Debug(yql_response)

	return TpbGetResults(yql_response, yql_format)
}

var (
	yql     = flag.String("yql", "https://query.yahooapis.com/v1/public/yql?q=", "yahoo query location")
	yql_qry = flag.String("yql_qry", "SELECT * FROM html WHERE url=<url> AND xpath='//tr'", "yahoo query string")

	tpb = flag.String("tpb", "https://thepiratebay.se/", "tpb web location")
)

func TpbCreateRequestUrlUsingYql(criteria TpbSearchCriteria, yql_format string) string {
	tpb_query := "\"" + *tpb + criteria.Method + "/" + criteria.Name + "/" + criteria.Page + "/" + criteria.Order + "/" + criteria.Category + "/" + "\""
	tpb_query = strings.Replace(tpb_query, " ", "%20", -1)
	yql_query := strings.Replace(*yql_qry, "<url>", tpb_query, 1)
	yql_query_escaped_plus := url.QueryEscape(yql_query)
	yql_query_escaped := strings.Replace(yql_query_escaped_plus, "+", "%20", -1)
	yql_str := *yql + yql_query_escaped + yql_format

	log.Debug("[CreateRequestUrlTpbYql] [yql_query_escaped]\n%s", tpb_query)

	return yql_str
}
