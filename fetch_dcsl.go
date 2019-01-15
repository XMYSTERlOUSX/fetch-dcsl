package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/bertrandmt-nflx/fetch-dcsl/widevine"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func todays_dcsl_file_name() string {
	return fmt.Sprintf("%04d%02d%02d-dcsl.bin", time.Now().Year(), time.Now().Month(), time.Now().Day())
}

// fetch a fresh DCSL from its server
func fetch_dcsl_data() []byte {
	const TheUrl = "https://www.googleapis.com/certificateprovisioning/v1/devicecertificatestatus/list?key=AIzaSyDMLcE1tgmHw8Eg5rUvrdPFgXT6VQl-rHQ"

	resp, err := http.Post(TheUrl, "application/x-www-form-urlencoded", &bytes.Buffer{})
	check(err)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	check(err)

	var dcsl_response_body map[string]string
	err = json.Unmarshal(data, &dcsl_response_body)
	check(err)

	signedList, prs := dcsl_response_body["signedList"]
	if !prs {
		panic(errors.New("no \"signedList\" element in JSON response"))
	}

	decoded, err := base64.URLEncoding.DecodeString(signedList)
	check(err)

	// save as "dcsl.bin" in current directory
	err = ioutil.WriteFile(todays_dcsl_file_name(), decoded, 0644)
	check(err)

	return decoded
}

func read_dcsl_data(r io.Reader) []byte {
	data, err := ioutil.ReadAll(r)
	check(err)
	return data
}

func main() {
	var pFetch = flag.Bool("f", false, "fetch fresh DCSL")
	var pSystemId = flag.Int("s", 0, "report on specific system ID")
	var pListManufacturers = flag.Bool("m", false, "list all manufacturers as CSV")
	flag.Parse()

	var data []byte
	if r, err := os.Open(todays_dcsl_file_name()); err == nil {
		data = read_dcsl_data(r)
	} else if *pFetch {
		data = fetch_dcsl_data()
	} else {
		data = read_dcsl_data(os.Stdin)
	}

	sdcsl := &widevine.SignedDeviceCertificateStatusList{}
	err := proto.Unmarshal(data, sdcsl)
	check(err)

	if 0 != *pSystemId {
		dcs_array := sdcsl.GetCertificateStatusList().GetCertificateStatus()
		found := false
		for _, dcs := range dcs_array {
			if dcs.GetDeviceInfo().GetSystemId() == uint32(*pSystemId) {
				found = true
				m := &jsonpb.Marshaler{}
				json, err := m.MarshalToString(dcs)
				check(err)

				fmt.Println(json)
				return
			}
		}
		if !found {
			panic(errors.New(fmt.Sprintln("can't find device certificate entry with system ID", *pSystemId)))
		}
		return
	}

	if *pListManufacturers {
		dcs_array := sdcsl.GetCertificateStatusList().GetCertificateStatus()
		manufacturers := make(map[string]int)
		for _, dcs := range dcs_array {
			manufacturer := dcs.GetDeviceInfo().GetManufacturer()
			if "" != manufacturer {
				manufacturers[manufacturer]++
			}
		}
		for manufacturer, device_count := range manufacturers {
			fmt.Printf("%s, %d\n", manufacturer, device_count)
		}
		return
	}

	m := &jsonpb.Marshaler{}
	json, err := m.MarshalToString(sdcsl)
	check(err)

	fmt.Println(json)
	return
}
