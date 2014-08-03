package netgo

import (
	"fmt"
)


func WTHisThisJSON(f interface{}) {
    switch vf := f.(type) {
    case map[string]interface{}:
        fmt.Println("is a map:")
        for k, v := range vf {
            switch vv := v.(type) {
            case string:
                fmt.Printf("%v: is string - %q\n", k, vv)
            case int:
                fmt.Printf("%v: is int - %q\n", k, vv)
            default:
                fmt.Printf("%v: ", k)
                WTHisThisJSON(v)
            }
        }
		
    case []interface{}:
        fmt.Println("is an array:")
        for k, v := range vf {
            switch vv := v.(type) {
            case string:
                fmt.Printf("%v: is string - %q\n", k, vv)
            case int:
                fmt.Printf("%v: is int - %q\n", k, vv)
            default:
                fmt.Printf("%v: ", k)
                WTHisThisJSON(v)
            }
        }
    }
}


