package netgo

import (
)



type SearchCriteria struct {
	Name string
	Method MethodType
	Category CategoryType
	Order OrderType
	Page string
}

type MethodType string
const (
	Method_TpbFile			MethodType = "TpbFile"
	Method_TpbUser			MethodType = "TpbUser"
)
var MethodList = [...]MethodType {
	Method_TpbFile,
	Method_TpbUser,
}


type CategoryType string
const (
	Category_Uncategorized		CategoryType = "Uncategorized"
	Category_Audio				CategoryType = "Audio"
	Category_Audio_Other		CategoryType = "Audio_Other"
	Category_Video				CategoryType = "Video"
	Category_Video_Other		CategoryType = "Video_Other"
	Category_Apps				CategoryType = "Apps"
	Category_Games				CategoryType = "Games"
	Category_Nsfw				CategoryType = "Nsfw"
	Category_Other				CategoryType = "Other"
)
var CategoryList = [...]CategoryType {
	Category_Uncategorized,
	Category_Audio,
	Category_Audio_Other,
	Category_Video,
	Category_Video_Other,
	Category_Apps,
	Category_Games,
	Category_Nsfw,
	Category_Other,
}


type OrderType string
const (
	Order_None			OrderType = "None"
	Order_Name_A		OrderType = "Name_A"
	Order_Name_Z		OrderType = "Name_Z"
	Order_Date_New		OrderType = "Date_New"
	Order_Date_Old		OrderType = "Date_Old"
	Order_Size_Big		OrderType = "Size_Big"
	Order_Size_Small	OrderType = "Size_Small"
	Order_Seed_Most		OrderType = "Seed_Most"
	Order_Seed_Least	OrderType = "Seed_Least"
	Order_Leech_Most	OrderType = "Leech_Most"
	Order_Leech_Least	OrderType = "Leech_Least"
	Order_Category_A	OrderType = "Category_A"
	Order_Category_Z	OrderType = "Category_Z"
	Order_Unordered		OrderType = "Unordered"
)
var OrderList = [...]OrderType {
	Order_None,
	Order_Name_A,
	Order_Name_Z,
	Order_Date_New,
	Order_Date_Old,
	Order_Size_Big,
	Order_Size_Small,
	Order_Seed_Most,
	Order_Seed_Least,
	Order_Leech_Most,
	Order_Leech_Least,
	Order_Category_A,
	Order_Category_Z,
	Order_Unordered,
}



func GetMethod(s string) MethodType {
	switch s {
	case string(Method_TpbFile): return Method_TpbFile
	case string(Method_TpbUser): return Method_TpbUser
	}
	return Method_TpbFile
}
func GetCategory(s string) CategoryType {
	switch s {
	case string(Category_Uncategorized): return Category_Uncategorized
	case string(Category_Audio): return Category_Audio
	case string(Category_Audio_Other): return Category_Audio_Other
	case string(Category_Video): return Category_Video
	case string(Category_Video_Other): return Category_Video_Other
	case string(Category_Apps): return Category_Apps
	case string(Category_Games): return Category_Games
	case string(Category_Nsfw): return Category_Nsfw
	case string(Category_Other): return Category_Other
	}
	return Category_Uncategorized
}
func GetOrder(s string) OrderType {
	switch s {
	case string(Order_None): return Order_None
	case string(Order_Name_A): return Order_Name_A
	case string(Order_Name_Z): return Order_Name_Z
	case string(Order_Date_New): return Order_Date_New
	case string(Order_Date_Old): return Order_Date_Old
	case string(Order_Size_Big): return Order_Size_Big
	case string(Order_Size_Small): return Order_Size_Small
	case string(Order_Seed_Most): return Order_Seed_Most
	case string(Order_Seed_Least): return Order_Seed_Least
	case string(Order_Leech_Most): return Order_Leech_Most
	case string(Order_Leech_Least): return Order_Leech_Least
	case string(Order_Category_A): return Order_Category_A
	case string(Order_Category_Z): return Order_Category_Z
	case string(Order_Unordered): return Order_Unordered
	}
	return Order_Unordered
}
/*


type ECategory int
const (
	_ ECategory = iota
	ECatUncategorized
	ECatAudio
	ECatVideo
	ECatApp
	ECatGame
	ECatPo
	ECatOther
)
var CategoryNames = [...]string {
	ECatUncategorized	: "Uncategorized",
	ECatAudio			: "Audio",
	ECatVideo			: "Video",
	ECatApp				: "App",
	ECatGame			: "Game",
	ECatPo				: "Po",
	ECatOther			: "Other",
}
*/


type SearchResult struct {
	Category string
	Name string
	Link string
	Magnet string
	Torrent string
	Comments int
	Seed int
	Leech int
	Trusted bool
	VIP bool
	CoverImage bool
	User string
	Date string
	Size string
}

