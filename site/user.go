package site

// Player Hello
type Player struct {
	FacebookID  string
	Name        string
	Summoners   []Summoner
	Permissions map[string]bool
}

type Summoner struct {
	ID  int64
	Ign string
}

func (u *Player) HasPermission(name string) bool {
	_, ok := u.Permissions[name]
	return ok
}

func (u *Player) AddPermission(name string) {
	u.Permissions[name] = true
}

func (u *Player) AddSummoner(account Summoner) {
	u.Summoners = append(u.Summoners, account)
}

func (u *Player) AddSummonerFromInfo(id int64, ign string) {
	summoner := Summoner{ID: id, Ign: ign}
	u.AddSummoner(summoner)
}
