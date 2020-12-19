package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urlList = []string{
	"http://api.turinglabs.net/api/v1/jd/jxfactory/create/zFFlhAJjtMI2V831ukVbeA==/",
	"http://api.turinglabs.net/api/v1/jd/jxfactory/create/yq67HFJ4c31f4le6SxIQyw==/",
	"http://api.turinglabs.net/api/v1/jd/jxfactory/count/",


	"http://api.turinglabs.net/api/v1/jd/bean/create/fw4h7b36sezhaflagqig2l6nea/",
	"http://api.turinglabs.net/api/v1/jd/bean/create/3hffvb4ywqkfk6fmrmnuugt3ne/",
	"http://api.turinglabs.net/api/v1/jd/bean/count/",

	"http://api.turinglabs.net/api/v1/jd/farm/create/c56897a51a194b3682e668e1256a349d/",
	"http://api.turinglabs.net/api/v1/jd/farm/create/cf43943800b94751ae63f14a3f3a3938/",
	"http://api.turinglabs.net/api/v1/jd/farm/count/",

	"http://api.turinglabs.net/api/v1/jd/pet/create/MTAxODExNTM5NDAwMDAwMDAzOTc0OTkyNw==/",
	"http://api.turinglabs.net/api/v1/jd/pet/create/MTE1NDQ5OTIwMDAwMDAwMzk0NzUwNjk=/",
	"http://api.turinglabs.net/api/v1/jd/pet/count/",

	"http://api.turinglabs.net/api/v1/jd/ddfactory/create/P04z54XCjVWnYaS5jAKC2n63HlPlA/",
	"http://api.turinglabs.net/api/v1/jd/ddfactory/create/P04z54XCjVWnYaS5nxBXCmkgCQ/",
	"http://api.turinglabs.net/api/v1/jd/ddfactory/count/",

	"http://api.turinglabs.net/api/v1/jd/cleantimeinfo/",
}

func main() {
	for _, api := range urlList {
		time.Sleep(time.Second * 10)
		_ = withRetry(50, time.Millisecond*100, func() error {
			return sendGetRequest(api)
		})
	}
}

type Response struct {
	Code        int         `json:"code"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
	PoweredBy   string      `json:"powered by"`
	SponsoredBy string      `json:"sponsored by"`
}

func sendGetRequest(api string) error {
	resp, err := http.Get(api)
	if err != nil {
		fmt.Printf("failed to send api reqeust, url=%s, err=%v\n", api, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("http code error, api=%s, code=%d, msg=%s", api, resp.StatusCode, resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to get api resp body, url=%s, err=%v\n", api, err)
		return err
	}

	fmt.Printf("api: %s, resp body: %s\n", api, string(data))

	rspBody := &Response{}
	err = json.Unmarshal(data, rspBody)
	if err != nil {
		fmt.Printf("failed to unmarshal resp body, url=%s, err=%v\n", api, err)
		return err
	}

	if rspBody.Code != 200 {
		if rspBody.Code == 400 && strings.Contains(rspBody.Message, "existed") {
			return nil
		}

		return fmt.Errorf("failed to send request, url=%s, err=%s", api, fmt.Sprintf("code=%d, err=%s", rspBody.Code, rspBody.Message))
	}

	return nil
}

func withRetry(attempts uint, initSleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if attempts--; attempts > 0 {
			jitter := time.Duration(rand.Int63n(int64(initSleep)))
			initSleep += jitter / 2

			time.Sleep(initSleep)
			return withRetry(attempts, 2*initSleep, f)
		}

		return err
	}

	return nil
}
