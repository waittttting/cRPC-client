package conf

type LocalConf struct {
	Client Client
}

type Client struct {
	ServerName 			string
	ServerVersion 		string
	ConfigCenterHost 	string
}
