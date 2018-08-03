package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var login = ""
var password = ""
var token = ""

type auth struct {
	Email           string `json:"email"`
	ID              int    `json:"id"`
	LastAppActivity string `json:"lastAppActivity"`
	LastKnownArea   string `json:"lastKnownArea"`
	Name            string `json:"name"`
}

type activity struct {
	PetActivityAim   int    `json:"petActivityAim"`
	PetAvatarPath    string `json:"petAvatarPath"`
	BatteryCharge    int    `json:"batteryCharge"`
	SubscriptionDate string `json:"subscriptionDate"`
	CurrentEnergy    int    `json:"currentEnergy"`
	CurrentActivity  int    `json:"currentActivity"`
	CurrentDistance  int    `json:"currentDistance"`
	UpdateDate       string `json:"updateDate"`
	DeviceStatus     string `json:"deviceStatus"`
	LightStatus      bool   `json:"lightStatus"`
	PetID            int    `json:"petId"`
	SignalLevel      string `json:"signalLevel"`
}

type pet struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Birthday    int     `json:"birthday"`
	Weight      float32 `json:"weight"`
	Height      float32 `json:"height"`
	ActivityAim int     `json:"activityAim"`
}

func main() {
	flag.StringVar(&login, "login", "", "Mishiko login/email")
	flag.StringVar(&password, "password", "", "Mishiko password")
	flag.StringVar(&token, "token", "", "Mishiko token")

	flag.Parse()

	if token == "" && login == "" && password == "" {
		log.Println("You should provide token or login/password to start")
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			pets := getPets(false)
			if len(pets) > 0 {
				for index := range pets {
					petData := pets[index]
					petActivity := getActivity(petData.ID, false)
					log.Println(petActivity.PetID, petActivity.CurrentEnergy, petActivity.CurrentActivity, petActivity.CurrentDistance)
				}
			}
		})

		http.ListenAndServe(":8081", nil)
	}
}

func doAuth(login string, password string) (authtoken string, err error) {
	if login != "" && password != "" {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://api2.mishiko.intech-global.com/profile/auth?email="+login+"&pass="+password+"&type=IOS", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-SPOTTY-AUTH-NEW", "X-SPOTTY-AUTH-NEW")
		resp, err := client.Do(req)
		defer resp.Body.Close()

		if err != nil {
			return "", err
		}

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
		} else {
			log.Println("got token", authtoken)
		}
	} else {
		// log.Println
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api2.mishiko.intech-global.com/devpet/list?timezone=3", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SPOTTY-AUTH-NEW", "X-SPOTTY-AUTH-NEW")
	req.Header.Set("X-SPOTTY-ACCESS-TOKEN", token)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
		return
	}
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
		} else {
			log.Println("got token", authtoken)
		}
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api2.mishiko.intech-global.com/devpet/"+strconv.Itoa(petID)+"/main_data?timezone=3", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SPOTTY-AUTH-NEW", "X-SPOTTY-AUTH-NEW")
	req.Header.Set("X-SPOTTY-ACCESS-TOKEN", token)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
		return
	}
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
	}
	return
}
