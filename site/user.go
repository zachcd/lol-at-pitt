package site

type User struct {
	FacebookId  string
	LeagueId    int64
	Name        string
	Permissions map[string]bool
}

type Permission struct {
	Name string
}

func (u *User) HasPermission(name string) bool {
	_, ok := u.Permissions[name]
	return ok
}

func (u *User) IsCaptain() bool {
	return u.HasPermission("captain")
}
