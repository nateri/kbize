package netgo

import (
	//"fmt"
	//"os"
	"net/url"
)


type CriteriaFromQuery struct {
	Name     string
	Method   string
	Category string
	Order    string
	Page     string
	Debug    string
}

func GetSearchCriteria(query CriteriaFromQuery) SearchCriteria {
	var criteria SearchCriteria
	
	criteria.Name = query.Name
	criteria.Method = Method_TpbFile
	criteria.Category = GetCategory(query.Category)
	criteria.Order = GetOrder(query.Order)
	criteria.Page = query.Page
	
	return criteria
}

func GetTpbSearchCriteriaFromUrl(u *url.URL) (TpbSearchCriteria, string) {
	dbg := ""

	url_params := GetParamsFromUrl(u)
	search_criteria := GetSearchCriteria(url_params)
	tpb_criteria := ConvertCriteriaForTpb(search_criteria)
	
	return tpb_criteria, dbg
}

func GetParamsFromUrl(u *url.URL) (CriteriaFromQuery) {
	var criteria CriteriaFromQuery
	criteria.Name = ""
	criteria.Category = "Uncategorized"
	criteria.Order = ""
	criteria.Method = "TpbFile"
	criteria.Page = "0"
	criteria.Debug = ""

	if u == nil {
		return criteria
	}

	m, _ := url.ParseQuery(u.RawQuery)

	name, valid := m["search"]
	if valid && len(name) > 0 {
		criteria.Name = name[0]
	}
	method, valid := m["method"]
	if valid && len(method) > 0 {
		criteria.Method = method[0]
	}
	category, valid := m["category"]
	if valid && len(category) > 0 {
		criteria.Category = category[0]
	}
	order, valid := m["order"]
	if valid && len(order) > 0 {
		criteria.Order = order[0]
	}
	page, valid := m["page"]
	if valid && len(page) > 0 {
		criteria.Page = page[0]
	}
	test, valid := m["test"]
	if valid && len(test) > 0 {
		criteria.Debug = test[0]
	}

	return criteria
}


