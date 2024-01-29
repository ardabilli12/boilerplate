package healthCheck

type HttpResponseHealthCheck struct {
	RedisAvailable   bool `json:"redisAvailable"`
	MysqlAvailable   bool `json:"mysqlAvailable"`
	ServiceAvailable bool `json:"serviceAvailable"`
}
