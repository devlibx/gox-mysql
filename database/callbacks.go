package database

import "context"

type PostCallbackData struct {
	Ctx       context.Context `json:"-"`
	Name      string          `json:"name"`
	StartTime int64           `json:"start_time"`
	EndTime   int64           `json:"end_time"`
	TimeTaken int64           `json:"time_taken"`
}

func (p *PostCallbackData) GetDbCallNameForTracing() string {
	if p.Ctx != nil && p.Ctx.Value("__SQLCX_DB_CALL_NAME__") != nil {
		if val, ok := p.Ctx.Value("__SQLCX_DB_CALL_NAME__").(string); ok {
			return "Slow_Query_Trace__" + val
		}
	}
	return "Slow_Query_Trace__" + p.Name
}

type PostCallbackFunc func(data PostCallbackData)
