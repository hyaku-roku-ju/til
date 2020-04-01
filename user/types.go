package user

type PreferredTime struct {
	Hour int
	Min  int
}

type User struct {
	Id            string
	PreferredTime PreferredTime
}
