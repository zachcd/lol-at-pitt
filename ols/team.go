package ols

type Team struct {
	Name    string
	Players []int64
	Captain int64
	Wins    int
	Losses  int
	League  string
}

type Teams []*Team

// Sorting functions
func (p Teams) Len() int {
	return len(p)
}

func (p Teams) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Teams) Less(i, j int) bool {
	if p[i].Wins == p[j].Wins {
		return p[i].Losses < p[j].Losses
	} else {
		return p[i].Wins > p[j].Wins
	}
}

func (t *Team) IsPlayerOnTeam(playerId int64) bool {
	for _, id := range t.Players {
		if id == playerId {
			return true
		}
	}

	return false
}
