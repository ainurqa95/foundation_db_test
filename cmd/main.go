package main

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/google/uuid"

	"log"
	"net/http"
)

type HttpServer struct {
	fdb fdb.Database
}

type Wallet struct {
	Uid      string
	Sum      int64
	Address  string
	MemberId string
	Payload  string
}

func (s *HttpServer) insertFoundation(w http.ResponseWriter, r *http.Request) {
	_, _ = s.fdb.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
		firstUid := uuid.New().String()
		secondUid := uuid.New().String()
		firstWallet := Wallet{Uid: firstUid, Sum: 123, MemberId: firstUid, Payload: "SOME {{repeat(5, 7)}}',\\\\n  {\\\\n    _id: '{{objectId()}}',\\\\n    index: '{{index()}}',\\\\n    guid: '{{guid()}}',\\\\n    isActive: '{{bool()}}',\\\\n    balance: '{{floating(1000, 4000 "}
		secondWallet := Wallet{Uid: secondUid, Sum: 125, MemberId: firstUid, Payload: "SOME {{repeat(5, 7)}}',\\\\n  {\\\\n    _id: '{{objectId()}}',\\\\n    index: '{{index()}}',\\\\n    guid: '{{guid()}}',\\\\n    isActive: '{{bool()}}',\\\\n    balance: '{{floating(1000, 4000 "}
		tr.Set(fdb.Key("id"+firstUid), []byte(fmt.Sprintf("%v", firstWallet)))
		tr.Set(fdb.Key("id"+secondUid), []byte(fmt.Sprintf("%v", secondWallet)))

		return
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {
	// Different API versions may expose different runtime behaviors.
	fdb.MustAPIVersion(630)

	foundb, err := fdb.OpenDatabase("/etc/foundationdb/fdb.cluster")
	if err != nil {
		log.Println(err)
		return
	}
	srv := HttpServer{
		fdb: foundb,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/founddb", srv.insertFoundation)

	go func() {
		err := http.ListenAndServe(":3334", mux)
		log.Println(err)
	}()

	select {}
}
