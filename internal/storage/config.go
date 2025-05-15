package postgres

type DatabaseConfig struct {
	Username string `json:"username"`
	DBName   string `json:"dbname"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}
