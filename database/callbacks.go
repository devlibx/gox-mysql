package database

type PostCallbackData struct {
	Name      string `json:"name"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	TimeTaken int64  `json:"time_taken"`
}

type PostCallbackFunc func(data PostCallbackData)
