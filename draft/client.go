//Draft client processes command object, stores it, and runs the command. The basic commands are: bid, win, undo win, next auction, undo next auction (previous auction)

//should handle storage, and results. all the DraftManager (or whatever you deem to name it) should do is handle Auctioners, bid-items and the ability to reverse/move forward with those decisions i.e. winning an auction subtracts points undoing an auction adds points etc. not implementing "bids on this player" since bids are temporal- bids are temporary, winning is forever!
package draft

type Command struct {
	command_name string
	amount       int
	curr_captain DraftPlayer
}

/*
//accessor methods
func GetCurrPlayer() (DraftPlayers, DraftPlayer) {
	var players DraftPlayers
	players = GetPlayers() //returns a list of ols.Player
	curr_player := players[0]
	return players, *curr_player
}

func GetCurrAuctioners() DraftCaptains {
	var captains DraftCaptains
	captains = GetCaptains() //returns a list of ols.Player
	return captains
}

func InitNewDraft() (chan command, chan command) {

	captains := GetCurrAuctioners()
	players, curr_player := GetCurrPlayer()

	command_in := make(chan<- Command)
	results_out := make(<-chan Command)

	go processing()
	return command_in, results_out
}

//processes commands
func processing() {
	log.Printf("beginning processing...")
	for {
		var received_command Command
		var highest_bidder DraftCaptain
		var highest_bid int
		var curr_player DraftPlayer
		/*
			received_command = <-command_in
			// TODO store received_command

			if received_command.command_name == "bid" {
				highest_bidder, highest_bid = bid(received_command.curr_captain, received_command.amount, highest_bidder, highest_bid)
			} else if received_command.command_name == "win" {
				highest_bidder = win(highest_bidder, highest_bid)
			} else if received_command.command_name == "undo_win" {
				undo_win()
			} else if received_command.command_name == "next_auction" {
				curr_player = next_auction()
			} else if received_command.command_name == "undo_next_auction" {
				curr_player = next_auction()
			} else {
				log.Printf("Error: unknown command")
			}
			results_out <- highest_bidder
			results_out <- highest_bid
			results_out <- curr_player
	}
}

//Commands

//bid on player with X amount
func bid(curr_bidder string, X int, highest_bidder string, highest_bid int) (string, int) {
	if X > curr_bidder.Points {
		log.Printf("Error: Not enough points to place this bid")
	} else if X > highest_bid {
		highest_bid = X
		highest_bidder = curr_bidder
	}
	return highest_bidder, highest_bid
}

//end current player's draft, takes the top bidder and removes that many points from his current pool of money, and adds that player to that captain's roster
func win(highest_bidder DraftPlayer, highest_bid int) string {
	//end current player's draft
	winner := highest_bidder
	top_bid := highest_bid
	//remove X points from top bidder's money pool
	winner.Points = winner.Points - top_bid
	//add current player to captain's roster
	curr_player.HighestBidder = winner.Team
	//remove current player from list of upcoming players
	players.Delete(curr_player)
	return highest_bidder
}

//undo the actions of win
func undo_win() {
	//add X points to top bidder's money pool
	winner.Points = winner.Points + curr_bid
	//remove current player from captain's roster
	curr_player.HighestBidder = nil
	//TODO add player back to list of players
}

//start the next player's draft
func next_auction() DraftPlayer {
	_, curr_player = GetCurrPlayer()
	return curr_player
}

//goes back to previous player's draft
func undo_next_auction() DraftPlayer {
	//_, curr_player = GetCurrPlayer()
	return curr_player
}
*/
