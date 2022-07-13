package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

func CheckTaskStatus(taskId int, debug bool) *TaskCheckResponse {
	var response *fasthttp.Response
	var body []byte
	var err error

	var request = fasthttp.AcquireRequest()

	body, err = json.Marshal(&CheckTaskStatusRequest{
		TaskId:  taskId,
		License: Config.License,
	})
	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred preparing request to check task status. Error: %s\n",
			err.Error())
		return nil
	}

	request.SetRequestURI(GetURL("/tasks/check_task"))
	request.Header.SetMethod("POST")
	request.Header.SetContentType("application/json")
	request.SetBody(body)

	response, err = SendRequest(request)
	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred sending request to check task status. Error: %s\n",
			err.Error())
		return nil
	}

	defer fasthttp.ReleaseResponse(response)

	var taskResult = TaskCheckResponse{}

	if debug {
		fmt.Println(string(response.Body()))
	}

	err = json.Unmarshal(response.Body(), &taskResult)
	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred parsing response to check task status. Error: %s\n",
			err.Error())
		return nil
	}

	return &taskResult
}

type Account struct {
	Name     string
	Password string
	Claimed  bool
}

func (A *Account) ClaimUsername(username string) {
	var err error
	var body []byte
	var response *fasthttp.Response

	var request = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	request.Header.SetMethod("POST")
	request.SetRequestURI(GetURL("/tasks/claim_username"))

	body, err = json.Marshal(&UsernameTaskRequest{
		Username:    username,
		AccountName: A.Name,
		License:     Config.License,
	})
	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred preparing request to claim username: %s. Error: %s\n",
			username, err.Error())
		return
	}

	request.Header.SetContentType("application/json")
	request.SetBody(body)

	response, err = SendRequest(request)
	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred sending request to claim username: %s. Error: %s\n",
			username, err.Error())
		return
	}

	defer fasthttp.ReleaseResponse(response)

	var taskResponse = &NewTaskResponse{}

	err = json.Unmarshal(response.Body(), taskResponse)
	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred parsing response to claim username: %s. Error: %s\n",
			username, err.Error())
		return
	}

	var taskResult = &TaskCheckResponse{
		Status: "Wait",
	}

	var counter = 0

	for taskResult.Status == "Wait" {
		taskResult = CheckTaskStatus(taskResponse.TaskId, true)

		if taskResult == nil {
			return
		}

		if taskResult.Status == "Error" {
			var errMsg string
			var ok bool

			errMsg, ok = taskResult.Response.(map[string]interface{})["Response"].(string)
			if !ok {

				var firstResponse = taskResult.Response
				if firstResponse != nil {
					var secondResponse = firstResponse.(map[string]interface{})["Response"]
					if secondResponse != nil {
						var dataResponse = secondResponse.(map[string]interface{})["data"]
						if dataResponse != nil {
							errMsg = dataResponse.(map[string]interface{})["description"].(string)
						}
					}
				}
				if !ok {
					errMsg = fmt.Sprintf("Unknown error. Response: %s", response.Body())
				}
			}

			log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred while claiming username: %s. Error: %s. Response: %s\n",
				username, errMsg, string(response.Body()))
		}

		if taskResult.Status == "Complete" {
			log.Printf("[\x1b[96mINFO\x1b[97m] Successfully claimed username: %s\n", username)
			if strings.Contains(string(response.Body()), username) {
				A.Claimed = true
				LogChannel <- fmt.Sprintf("Successfully claimed username: %s\n", username)
			}
		}

		if counter >= 3000 {
			log.Printf("[\x1b[91mERROR\x1b[97m] Failed to claim username after trying for 5 minutes.\n")
			return
		}

		counter++
	}
}

func (A *Account) CheckUsernameAvailable(request *fasthttp.Request, username string) bool {
	var response *fasthttp.Response
	var err error

	response, err = SendRequest(request)
	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred sending request to check username: %s. Error: %s\n",
			username, err.Error())
		return false
	}

	defer fasthttp.ReleaseResponse(response)

	var taskResponse = &NewTaskResponse{}

	err = json.Unmarshal(response.Body(), taskResponse)
	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred parsing response to check username: %s. Error: %s\n",
			username, err.Error())
		return false
	}

	var taskResult = &TaskCheckResponse{
		Status: "Wait",
	}

	var counter = 0

	for taskResult.Status == "Wait" {
		taskResult = CheckTaskStatus(taskResponse.TaskId, false)
		if taskResult == nil {
			return false
		}

		if taskResult.Status == "Error" {
			var errMsg string
			var ok bool

			errMsg, ok = taskResult.Response.(map[string]interface{})["Response"].(string)
			if !ok {

				var firstResponse = taskResult.Response
				if firstResponse != nil {
					var secondResponse = firstResponse.(map[string]interface{})["Response"]
					if secondResponse != nil {
						var dataResponse = secondResponse.(map[string]interface{})["data"]
						if dataResponse != nil {
							errMsg = dataResponse.(map[string]interface{})["description"].(string)
						}
					}
				}
				if !ok {
					errMsg = fmt.Sprintf("Unknown error. Response: %s", response.Body())
				}
			}

			log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred while checking username: %s. Error: %s. Response: %s\n",
				username, errMsg, string(response.Body()))
		}

		if taskResult.Status == "Complete" {
			var isFree, ok = taskResult.Response.(map[string]interface{})["Response"].(map[string]interface{})["is_valid"].(bool)
			if !ok {
				fmt.Println("Error occured in conversion")
			}
			fmt.Printf("%s isFree: %v\n", username, isFree)

			return true
		}

		if counter >= 3000 {
			log.Printf("[\x1b[91mERROR\x1b[97m] Failed to check username: %s after trying for 5 minutes. Check your proxy or contact support\n",
				username)
			return false
		}

		time.Sleep(100 * time.Millisecond)
		counter++
	}

	return false
}

func (A *Account) StartClaimer(username string) {
	var err error
	var body []byte

	var request = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	request.SetRequestURI(GetURL("/tasks/check_username"))

	body, err = json.Marshal(&UsernameTaskRequest{
		Username:    username,
		AccountName: A.Name,
		License:     Config.License,
	})

	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred preparing request to check username: %s. Error: %s\n",
			username, err.Error())
		return
	}

	request.Header.SetContentType("application/json")
	request.SetBody(body)
	request.Header.SetMethod("POST")

	for !A.Claimed {
		if A.CheckUsernameAvailable(request, username) {
			for i := 0; i < 5; i++ {
				A.ClaimUsername(username)
			}
			if !A.Claimed {
				log.Printf("[\x1b[96mINFO\x1b[97m] Failed to claim: %s.", username)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (A *Account) AuthenticateAccount() bool {
	var err error
	var body []byte
	var response *fasthttp.Response

	var request = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	request.SetRequestURI(GetURL("/tasks/login_tiktok"))
	request.Header.SetMethod("POST")

	log.Printf("[\x1b[96mINFO\x1b[97m] Logging in to account: %s:%s", A.Name, A.Password)

	body, err = json.Marshal(&TikTokLoginRequest{
		Username: A.Name,
		Password: A.Password,
		Proxy:    GetRandomProxy(),
		License:  Config.License,
	})
	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred preparing request for account: %s:%s. Error: %s\n",
			A.Name, A.Password, err.Error())
		return false
	}

	request.SetBody(body)
	request.Header.SetContentType("application/json")

	response, err = SendRequest(request)
	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] occurred sending request for account: %s:%s. Error: %s\n",
			A.Name, A.Password, err.Error())
		return false
	}

	defer fasthttp.ReleaseResponse(response)

	var taskResponse = &NewTaskResponse{}

	err = json.Unmarshal(response.Body(), taskResponse)

	if err != nil {
		log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred parsing response for account: %s:%s. Error: %s\n",
			A.Name, A.Password, err.Error())
		return false
	}

	var taskResult = &TaskCheckResponse{
		Status: "Wait",
	}

	var counter = 0

	for taskResult.Status == "Wait" {
		log.Printf("waiting for task...")
		taskResult = CheckTaskStatus(taskResponse.TaskId, false)
		if taskResult == nil {
			log.Printf("Invalid shit")
			return false
		}

		if taskResult.Status == "Error" {
			var errMsg string
			var ok bool

			errMsg, ok = taskResult.Response.(map[string]interface{})["Response"].(string)
			if !ok {

				var firstResponse = taskResult.Response
				if firstResponse != nil {
					var secondResponse = firstResponse.(map[string]interface{})["Response"]
					if secondResponse != nil {
						var dataResponse = secondResponse.(map[string]interface{})["data"]
						if dataResponse != nil {
							errMsg = dataResponse.(map[string]interface{})["description"].(string)
						}
					}
				}
				if !ok {
					errMsg = fmt.Sprintf("Unknown error")
				}
			}

			log.Printf("[\x1b[91mERROR\x1b[97m] Error occurred while authenticating account: %s:%s. Error: %s. Response: %s\n",
				A.Name, A.Password, errMsg, string(response.Body()))
		}

		if taskResult.Status == "Complete" {
			log.Printf("[\x1b[96mINFO\x1b[97m] Succcessfully logged in with acocunt %s:%s\n",
				A.Name, A.Password)
			return true
		}

		if counter >= 600 {
			log.Printf("[\x1b[91mERROR\x1b[97m] Failed to login after trying for 10 minutes.\n")
			return true
		}

		time.Sleep(10 * time.Second)
		counter++
	}

	return false
}

func (A *Account) Start() {
	A.Claimed = false

	ok := A.AuthenticateAccount()
	if !ok {
		ActiveSessionsWG.Done()
	}

	for _, username := range Usernames {
		go A.StartClaimer(username)
	}
}
