package config

//Config Configuration struct
type Config struct {
	URL    string `json:"url"`
	Port   int16  `json:"port"`
	DbURL  string `json:"dburl"`
	DbName string `json:"dbname"`
}
