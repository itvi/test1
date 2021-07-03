package main

import (
	"ams/pkg/auth"
	"ams/pkg/models"
	"ams/pkg/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Initialize add config information to database
func (app *application) InitializeComputerConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		app.render(w, r, "computer.config.init.html", &templateData{})
	}
	if r.Method == "POST" {
		jsonstring := r.PostFormValue("info")
		var computer []*models.ComputerConfig

		if err := json.Unmarshal([]byte(jsonstring), &computer); err != nil {
			app.errorLog.Println("json unmarshal error:", err)
			app.session.Put(r, "flash", "初始化失败:"+err.Error())
			return
		}

		if err := app.computer.Initialize(computer); err != nil {
			app.errorLog.Printf("Initialize(add config) error:%v", err)
			app.session.Put(r, "flash", "初始化失败:"+err.Error())
			return
		}
		w.Write([]byte("init OK"))
	}
}

func searchByIP(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	computerConfig, err := GetComputerConfig(ip)
	if err != nil {
		w.Write([]byte("查询失败！"))
		return
	}

	j, err := json.Marshal(computerConfig)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := struct {
		ComputerConfig *models.ComputerConfig
		GenericString  string
	}{
		computerConfig,
		string(j),
	}
	renderPartialView(w, data, "./ui/html/computer.config.html")
}

// search by ips with concurrence
func searchByIPs(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")

	var strJson string
	ipRange := strings.Split(ip, "-")
	// fmt.Println("IP range:", ipRange)

	ips := []string{}
	results := []models.ComputerConfig{}
	ch := make(chan models.ComputerConfig)

	if len(ipRange) == 2 {
		beginIP := util.IP2Int(ipRange[0])
		endIP := util.IP2Int(ipRange[1])
		if beginIP > endIP {
			log.Println("IP range error: begin > end")
		}

		for i := beginIP; i <= endIP; i++ {
			ip := util.Int2IP(i)
			ips = append(ips, ip)
		}

		for _, ip := range ips {
			go GetComputerConfigs(ip, ch)
		}

		for range ips {
			result := <-ch
			results = append(results, result)
		}
	}

	if len(ipRange) == 1 {
		go GetComputerConfigs(ip, ch)

		result := <-ch
		results = append(results, result)
	}

	// pass json to html page
	var j []byte
	var err error

	j, err = json.MarshalIndent(results, " ", "")

	if err != nil {
		log.Println("Marshall error:", err)
	}
	strJson = string(j)

	data := struct {
		ComputerConfig []models.ComputerConfig
		GenericString  string
	}{
		results, strJson,
	}

	var tplName = "./ui/html/computer.configs.html"
	renderPartialView(w, data, tplName)
}

func GetComputerConfig(ip string) (*models.ComputerConfig, error) {
	token, err := auth.GenerateToken()
	if err != nil {
		log.Println("Generate token error:", err.Error())
		return nil, err
	}
	client := &http.Client{}
	client.Timeout = time.Second * 10
	req, err := http.NewRequest("GET", "http://"+ip+":10000", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Token", token)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	info := &models.ComputerConfig{}
	if err := decoder.Decode(&info); err != nil {
		log.Println(err)
		return nil, err
	}

	info.IP = ip
	return info, nil
}

// query computer config information by token
func GetComputerConfigs(ip string, ch chan<- models.ComputerConfig) {
	token, err := auth.GenerateToken()
	if err != nil {
		fmt.Println("Generate token error:", err.Error())
		return
	}

	var config models.ComputerConfig

	client := &http.Client{}
	client.Timeout = time.Second * 10
	// request
	req, err := http.NewRequest("GET", "http://"+ip+":10000", nil)
	if err != nil {
		fmt.Printf("Request error: %v\n", err)
		ch <- config
	} else {
		req.Header.Set("Token", token)
		// response
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Response error: %v\n", err)
			config.IP = ip
			ch <- config
		} else {
			defer resp.Body.Close()
			decoder := json.NewDecoder(resp.Body)
			if err = decoder.Decode(&config); err != nil {
				fmt.Printf("Decoder decode error: %v\n", err)
			}
			config.IP = ip
			ch <- config
		}
	}
}
