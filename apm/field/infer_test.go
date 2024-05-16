package field

import (
	"encoding/json"
	"fmt"
	"testing"
)

type cust int

func TestNullPtr(t *testing.T) {
	var x *cust
	// fmt.Println(InferPrimitiveValueByReflect(reflect.ValueOf(x)))
	// fmt.Println(InferPrimitiveValueByReflect(reflect.Value{}))
	fmt.Println(Any("", x))

	// var y *int
	// // fmt.Println(InferPrimitiveValueByReflect(reflect.ValueOf(y)))
	// // fmt.Println(InferPrimitiveValueByReflect(reflect.Value{}))
	// fmt.Println(InferPrimitiveValue(y))
}

// 2,400,000,000
// 9,007,199,254,740,991
// 9,007,199,254,740,991
// 281,474,976,710,656

func TestAny(t *testing.T) {
	data := []byte(`[{
        "alg": null,
        "fee": 8,
        "ftype": 0,
        "score": 100.0,
        "copyrightId": 10,
        "mvid": 0,
        "transNames": null,
        "no": 7,
        "commentThreadId": "R_SO_10",
        "publishTime": 0,
        "status": 0,
        "privilege": null,
        "djid": 0,
        "album": {
            "id": 10,
            "name": "专辑",
            "picUrl": "url",
            "tns": [],
            "pic_str": "10",
            "pic": 10
        },
        "artists": [
            {
                "id": 6068,
                "name": "杨坤",
                "tns": [],
                "alias": []
            },
            {
                "id": 1024308,
                "name": "张碧晨",
                "tns": [],
                "alias": []
            }
        ],
        "alias": [],
        "type": 0,
        "duration": 373318,
        "name": "标题 (Live版)",
        "id": 10
    }]`)

	var obj any
	if err := json.Unmarshal(data, &obj); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("obj: %v\n", obj)

	f := Any("", obj)

	fmt.Printf("f: %#v\n", f)

}
