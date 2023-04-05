package database

type PostCallbackData struct {
	Name      string
	StartTime int64
	EndTime   int64
	TimeTaken int64
}

type PostCallbackFunc func(data PostCallbackData)
