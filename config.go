package jupiter

type Config struct {
	ListenPort   int
	HTTPPort     int
	IndexFiles   []string
	DataLogFiles []string
}

func readConfig(f string) *Config {
	return nil
}
