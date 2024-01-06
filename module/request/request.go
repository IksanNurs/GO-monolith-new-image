package request

type Request struct {
	ID                int    `json:"id"`
	UserID            int    `json:"user_id"`
	Route             string `json:"route"`
	Method            string `json:"method"`
	IPAddress         string `json:"ip_address"`
	RemoteAddr        string `json:"remote_addr"`
	HTTPXForwardedFor string `json:"http_x_forwarded_for"`
	HTTPUserAgent     string `json:"http_user_agent"`
	CreatedAt         int64  `json:"created_at"`
}
