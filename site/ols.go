package site

import "errors"

type OlsSplit struct {
	Teams   []Team
	Matches []Match
}

func (ols *OlsSplit) GetTeam(name string) (*Team, error) {
	for _, team := range ols.Teams {
		if team.Name == name {
			return &team, nil
		}
	}

	return nil, errors.New("No team by that name")
}

func (ols *OlsSplit) GetMatch(redTeamName string, blueTeamName string) (*Match, error) {
	for _, match := range ols.Matches {
		if match.Blue.Name == blueTeamName || match.Red.Name == redTeamName {
			return &match, nil
		}
	}

	return nil, errors.New("Couldn't find match")
}

func (ols *OlsSplit) GetMatchesByTeamName(name string) []*Match {
	matches := []*Match{}
	for _, match := range ols.Matches {
		if match.Blue.Name == name || match.Red.Name == name {
			matches = append(matches, &match)
		}
	}

	return matches
}

type Team struct {
	Captain Player
	Players []Player
	Name    string
	Wins    int
	Losses  int
	League  string
}

type Match struct {
	Red    Team
	Blue   Team
	Winner Team
	ID     string
}
