package netgo

import (
	//"log"
	//"github.com/op/go-logging"

)
//var log = logging.MustGetLogger("netgo")


func TestTpbSearchCriteria(criteria SearchCriteria) TpbSearchCriteria {

	tpb_criteria := ConvertCriteriaForTpb(criteria)
	log.Debug("%v", tpb_criteria)
	
	return tpb_criteria
}

func TestTpbSearch(criteria TpbSearchCriteria) []SearchResult {

	entries := TpbSearch(criteria, "&format=json")
	log.Debug("%v", entries)
	
	return entries
}

func TestTpbYqlJson(in string) []SearchResult {

	entries := TpbGetResultsJson(in)
	
	return entries
}

func TestNetgo() {

	// Verify SearchResults from a static json string
	tpbJsonOut := TestTpbYqlJson(TestTpbJsonStr)
	for index, entry := range tpbJsonOut {
		log.Debug("[%d] %+v\n", index, entry)
	}
	
	// Verify TpbSearchCriteria from a SearchCriteria
	tpbCriteriaOut := TestTpbSearchCriteria( SearchCriteria {
		Name : "linux",
		Category : Category_Audio,
		Order : Order_Seed_Most,
		Page : "0",
	})
	log.Debug("%+v", tpbCriteriaOut)
	
	// Verify SearchResults from a TpbSearchCriteria
	tpbSearchOut := TestTpbSearch(tpbCriteriaOut)
	for index, entry := range tpbSearchOut {
		log.Debug("[%d] %+v\n", index, entry)
	}
}

