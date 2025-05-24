package woodis

type RedisDB struct {
	Master     *WooDis
	Id         int
	stringKeys map[string]string
}
