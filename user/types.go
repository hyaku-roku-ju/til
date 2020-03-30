package user

type PreferredTime struct {
	hour int
	min  int
}

type User struct {
	id            string
	preferredTime PreferredTime
}
