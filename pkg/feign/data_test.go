package feign

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestGetNew(t *testing.T) {
	http.HandleFunc("/one/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		t.Log("/one/", id)

		n := R[News]{
			Code:    200,
			Message: "查询成功",
			Data: News{
				Id:      1,
				Content: "hello world",
			},
		}

		marshal, err := json.Marshal(n)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Write(marshal)
	})

	go http.ListenAndServe(":18080", nil)

	time.Sleep(1 * time.Second)
	fClients.clientMap[WorkData] = DiscoverInstance{
		Scheme:      "http",
		Ip:          "127.0.0.1",
		Port:        18080,
		ServiceName: WorkData,
	}

	getNew, err := GetNew(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(getNew)
}
