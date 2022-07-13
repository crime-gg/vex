package main

type AuthenticationRequest struct {
	Username string `json:"username"`
	License  string `json:"license"`
	HWID     string `json:"hwid"`
}

type TikTokLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Proxy    string `json:"proxy"`
	License  string `json:"license"`
}

type CheckTaskStatusRequest struct {
	TaskId  int    `json:"task_id"`
	License string `json:"license"`
}

type UsernameTaskRequest struct {
	Username    string `json:"username"`
	AccountName string `json:"account_name"`
	License     string `json:"license"`
}

type NewTaskResponse struct {
	TaskId int    `json:"_TaskId"`
	Status string `json:"Status"`
}

type TaskCheckResponse struct {
	TaskId   int         `json:"_TaskId"`
	Status   string      `json:"Status"`   // Wait, Error, Complete
	Response interface{} `json:"Response"` // Convert this to TaskResponseObject after Status is no longer "Wait"
}

type TaskResponseObject struct {
	ErrorCode int         `json:"ErrorCode"`
	Response  interface{} `json:"Response"`
	Status    bool        `json:"Status"`
}

type CheckUsernameResponse struct {
	IsValid    bool   `json:"is_valid"`
	StatusMsg  string `json:"status_msg"`
	StatusCode int    `json:"status_code"`
}

type ClaimUsernameResponse struct {
	Data struct {
		LoginName string `json:"login_name"`
	}
	Message string `json:"message"`
}
