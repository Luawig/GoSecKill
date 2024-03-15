package main

import (
	"GoSecKill/common"
	"GoSecKill/pkg/util"
	"errors"
	"net/http"
	"strconv"
	"sync"
)

var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

var port = "8081"

var hashConsistent *common.Consistent

type AccessControl struct {
	// slice of host
	sourceArray map[int]interface{}
	// rw mutex
	sync.RWMutex
}

var accessControl = &AccessControl{sourceArray: make(map[int]interface{})}

func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.sourceArray[uid]
}

func (m *AccessControl) SetNewRecord(uid int) {
	m.Lock()
	m.sourceArray[uid] = uid
	m.Unlock()
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}
	if hostRequest == localHost {
		return m.GetDataFromMap(uid.Value)
	}
	return m.GetDataFromOtherMap(hostRequest, req)
}

func (m *AccessControl) GetDataFromMap(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)

	if data != nil {
		return true
	}
	return false
}

func (m *AccessControl) GetDataFromOtherMap(host string, request *http.Request) bool {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return false
	}
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return false
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/check", nil)
	if err != nil {
		return false
	}

	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	response, err := client.Do(req)
	if err != nil {
		return false
	}

	if response.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

func Check(rw http.ResponseWriter, req *http.Request) {
	// do something
}

func Auth(rw http.ResponseWriter, req *http.Request) error {
	err := CheckUserInfo(req)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserInfo(req *http.Request) error {
	uidCookie, err := req.Cookie("uid")
	if err != nil {
		return err
	}
	signCookie, err := req.Cookie("sign")
	if err != nil {
		return err
	}
	signByte, err := util.Decrypt([]byte(signCookie.Value))
	if err != nil {
		return err
	}
	if uidCookie.Value != string(signByte) {
		return errors.New("validate sign failed")
	}
	return nil
}

func main() {
	hashConsistent = common.NewConsistent()

	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	filter := common.NewFilter()

	filter.RegisterFilterUri("/check", Auth)

	http.HandleFunc("/check", filter.Handle(Check))

	_ = http.ListenAndServe(":8080", nil)
}
