package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func doAuth(login string, password string) (authtoken string, err error) {
	if login != "" && password != "" {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://api2.mishiko.intech-global.com/profile/auth?email="+login+"&pass="+password+"&type=IOS", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-SPOTTY-AUTH-NEW", "X-SPOTTY-AUTH-NEW")
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		authtoken = resp.Header.Get("X-Spotty-Access-Token")

		if resp.StatusCode != 200 || authtoken == "" {
			return "", nil
		}
		token = authtoken
	}

	return authtoken, nil
}

func getPets(reauth bool) (pets []pet) {
	if token == "" {
		authtoken, err := doAuth("", "")
		if err != nil || authtoken == "" {
			log.Println(err)
			return
		}
		log.Println("got token", authtoken)
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api2.mishiko.intech-global.com/devpet/list?timezone="+strconv.Itoa(timezone), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SPOTTY-AUTH-NEW", "X-SPOTTY-AUTH-NEW")
	req.Header.Set("X-SPOTTY-ACCESS-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 401 {
		if reauth {
			return
		}

		token, err = doAuth(login, password)
		if err != nil || token == "" {
			log.Println(err)
			return
		}

		return getPets(true)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &pets)
	if err != nil {
		log.Println(err, string(body))
	}
	return
}

func getActivity(petID int, reauth bool) (activityData activity) {
	if token == "" {
		authtoken, err := doAuth(login, password)
		if err != nil || authtoken == "" {
			log.Println(err)
			return
		}
		log.Println("got token", authtoken)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api2.mishiko.intech-global.com/devpet/"+strconv.Itoa(petID)+"/main_data?timezone="+strconv.Itoa(timezone), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SPOTTY-AUTH-NEW", "X-SPOTTY-AUTH-NEW")
	req.Header.Set("X-SPOTTY-ACCESS-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 401 {
		if reauth {
			return
		}

		token, err = doAuth(login, password)
		if err != nil || token == "" {
			log.Println(err)
			return
		}

		return getActivity(petID, true)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &activityData)
	if err != nil {
		log.Println(err, string(body))
	} else {
		if activityData.BatteryCharge < 0 && activityData.BatteryCharge != -100 {
			activityData.Charging = true
		}
	}
	return
}

func getPetsLocations(reauth bool) (locations locations) {
	if token == "" {
		authtoken, err := doAuth("", "")
		if err != nil || authtoken == "" {
			log.Println(err)
			return
		}
		log.Println("got token", authtoken)
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api2.mishiko.intech-global.com/devpet/locations", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SPOTTY-AUTH-NEW", "X-SPOTTY-AUTH-NEW")
	req.Header.Set("X-SPOTTY-ACCESS-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 401 {
		if reauth {
			return
		}

		token, err = doAuth(login, password)
		if err != nil || token == "" {
			log.Println(err)
			return
		}

		return getPetsLocations(true)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &locations)
	if err != nil {
		log.Println(err, string(body))
	}
	return
}
