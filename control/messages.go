package yeelight

type Error struct {
	Code    int
	Message string
}

type Result struct {
	ID     int
	Result []string
	Error  Error
}
