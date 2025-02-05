package htmlbuilder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/CCPupp/pupper.moe/internal/discord"
	"github.com/CCPupp/pupper.moe/internal/player"
)

func BuildHTMLHeader() string {
	var finalHeader string = `<!DOCTYPE html>
	<html>
	<title>State Ranking Leaderboard</title>
	<meta charset="UTF-8" />
	<link rel="preconnect" href="https://fonts.gstatic.com">
	<link href="https://fonts.googleapis.com/css2?family=Roboto&display=swap" rel="stylesheet"> 
	<link rel="icon" href="../web/media/favicon.png" type="image/x-icon"/>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<link rel="stylesheet" type="text/css" href="../web/css/main.css" />
	<link rel="stylesheet" type="text/css" href="../web/css/flexbox.css" />
	<link rel="stylesheet" type="text/css" href="../web/css/normalize.css" />
	<link rel="stylesheet" type="text/css" href="../web/css/playercards.css" />
	<meta property="og:type" content="website">
	<meta property="og:title" content="State Leaderboard" />
	<meta property="og:description" content="A leaderboard for osu players split into each state" />
	<meta property="og:url" content="https://pupper.moe" />
	<meta property="og:image" content="full thumbnail path" />
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
	<script src="https://code.jquery.com/jquery-3.1.1.min.js"></script>
	<script src="../web/scripts/main.js"></script>`
	return finalHeader
}

func BuildHTMLFooter() string {
	var finalFooter string = `
	<script src="../web/scripts/localstorage.js"></script>
	</html>`
	return finalFooter
}

func BuildHTMLNavbar() string {
	finalString := `<body>
    <div class="navbar">
        <a href="/">Home</a>
        <a href="/all">All Users / Discords</a>
        <a href="https://twitter.com/ponparpanpor">Contact</a>
        <a href="/login">Customize My Card</a>
    </div>
	<br>
	<br>
	`
	return finalString
}

func CreateAllHTML() string {
	users := player.GetUserJSON()

	finalString := BuildHTMLHeader()
	finalString += BuildHTMLNavbar()
	finalString += `
	<br>
	<br>
	<div class='flex-container black-font'>
	<ol>
	<b>Total Users: ` + strconv.Itoa(len(users.Users)) + `</b><br><br>`

	for i := len(users.Users) - 1; i >= 0; i-- {
		finalString += ("<li><div style='height: 40px;' class='flex-center'><a href='/states/" + users.Users[i].State + "' class='usercard black-font'>" + users.Users[i].Username + "</a>")
		finalString += ("<b>State: " + users.Users[i].State + "</b></div></li>")
	}

	// Open our jsonFile
	discordJsonFile, err := os.Open("web/data/discords.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer discordJsonFile.Close()

	discordByteValue, _ := ioutil.ReadAll(discordJsonFile)

	// we initialize our Players array
	var discords discord.Discords

	json.Unmarshal(discordByteValue, &discords)
	finalString += `</ol><ol>`

	for i := 0; i < len(discords.Discords); i++ {
		finalString += `<a class="black-font" href=` + discords.Discords[i].Link + `> ` + discords.Discords[i].State + `'s Discord Server </a><br><br>`
	}

	finalString += `</div></body>`
	finalString += BuildHTMLFooter()

	return finalString
}

func CreateStateHTML(state string) string {
	discords := discord.GetDiscordJSON()
	discordString := ""
	for i := 0; i < len(discords.Discords); i++ {
		if discords.Discords[i].State == state {
			discordString += `<a href="` + discords.Discords[i].Link + `"> Discord Server </a>`
		}
	}
	var finalString = BuildHTMLHeader()
	users := player.GetUserJSON()

	users = player.SortUsers(users)

	finalString += `<body>
    <div class="navbar">
        <a href="/">Home</a>
		<a>`
	finalString += state
	finalString += `</a>
	` + discordString + `
	<a href="/login">Customize My Card</a>
    </div>
	<br>
	<br>
	<p id="result"></p>
	<div class="playerlist">
	`
	rank := 0
	for i := 0; i < len(users.Users); i++ {
		if (users.Users[i].State == state) && (users.Users[i].Statistics.Global_rank != 0) {
			rank++
			finalString += CreateUser(users.Users[i])
		}
	}
	finalString += "</div>"

	finalString += `</div></body>`
	finalString += BuildHTMLFooter()

	return finalString
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 2, 64)
}

func GetBackgroundText(bg player.User) string {
	if bg.Background == "true" {
		return "On"
	}
	return "Off"
}

func CreateOptions(user player.User) string {
	finalString := (`<div class="player-container black-font">
						<div class="user-settings">
							<div class="settings-info">
								<p>Hello ` + user.Username + `! Here you can change how your player card appears on the state leaderboard.</p>
								<input type="hidden" id="userid" value="` + strconv.Itoa(user.ID) + `"/>
								<br>
								<select id="bg">
									<option value="` + user.Background + `" selected hidden>` + GetBackgroundText(user) + `</option>
									<option value="true">On</option>
									<option value="false">Off</option>
								</select>
								<label>Background Image On/Off</label>
								<br>
								<br>
								<select id="mode" disabled>
									<option value="` + user.Playmode + `" selected hidden>` + user.Playmode + `</option>
									<option value="osu">osu</option>
									<option value="mania">mania</option>
									<option value="taiko">taiko</option>
									<option value="fruits">fruits</option>
								</select>
								<label>Gamemode Preference</label>
								<br>
								<br>
								<select id="state">
									<option value="` + user.State + `" selected hidden>` + user.State + `</option>
									<option value="Alabama">Alabama</option>
									<option value="Alaska">Alaska</option>
									<option value="Arizona">Arizona</option>
									<option value="Arkansas">Arkansas</option>
									<option value="California">California</option>
									<option value="Colorado">Colorado</option>
									<option value="Connecticut">Connecticut</option>
									<option value="Delaware">Delaware</option>
									<option value="Florida">Florida</option>
									<option value="Georgia">Georgia</option>
									<option value="Hawaii">Hawaii</option>
									<option value="Idaho">Idaho</option>
									<option value="Illinois">Illinois</option>
									<option value="Indiana">Indiana</option>
									<option value="Iowa">Iowa</option>
									<option value="Kansas">Kansas</option>
									<option value="Kentucky">Kentucky</option>
									<option value="Louisiana">Louisiana</option>
									<option value="Maine">Maine</option>
									<option value="Maryland">Maryland</option>
									<option value="Massachusetts">Massachusetts</option>
									<option value="Michigan">Michigan</option>
									<option value="Minnesota">Minnesota</option>
									<option value="Mississippi">Mississippi</option>
									<option value="Missouri">Missouri</option>
									<option value="Montana">Montana</option>
									<option value="Nebraska">Nebraska</option>
									<option value="Nevada">Nevada</option>
									<option value="New Hampshire">New Hampshire</option>
									<option value="New Jersey">New Jersey</option>
									<option value="New Mexico">New Mexico</option>
									<option value="New York">New York</option>
									<option value="North Carolina">North Carolina</option>
									<option value="North Dakota">North Dakota</option>
									<option value="Ohio">Ohio</option>
									<option value="Oklahoma">Oklahoma</option>
									<option value="Oregon">Oregon</option>
									<option value="Pennsylvania">Pennsylvania</option>
									<option value="Rhode Island">Rhode Island</option>
									<option value="South Carolina">South Carolina</option>
									<option value="South Dakota">South Dakota</option>
									<option value="Tennessee">Tennessee</option>
									<option value="Texas">Texas</option>
									<option value="Utah">Utah</option>
									<option value="Vermont">Vermont</option>
									<option value="Virginia">Virginia</option>
									<option value="Washington">Washington</option>
									<option value="West Virginia">West Virginia</option>
									<option value="Wisconsin">Wisconsin</option>
									<option value="Wyoming">Wyoming</option>
								</select>	
								<label>State Selection</label>
								<br>
								<br>
							<form action="/login">
								<button id="update">Submit</button>
							    <input type="submit" value="Refresh" />
							</form>
							</div>
						</div>`)
	return finalString
}

func CreateUser(user player.User) string {
	finalString := (`<div class="players-container" id="response">
						<div class="player">
							<div class="player-preview">
								<h4>#` + strconv.Itoa((player.GetUserStateRank(user.ID, user.State))) + `</h4>
								<image class="playerpfp" href="https://osu.ppy.sh/users/` + strconv.Itoa(user.ID) + `" src="http://s.ppy.sh/a/` + strconv.Itoa(user.ID) + `"></image>
							</div>
							<div class="player-info" style="` + GetBackground(user) + `">
								<div class="progress-container">
									<span class="progress-text">
										<h5>Mode: ` + user.Playmode + `</h5>
										<h5>Level ` + strconv.Itoa(user.Statistics.Level.Current) + `.` + strconv.Itoa(user.Statistics.Level.Progress) + `</h5>
										<h5>Discord: ` + user.Discord + `</h5>
									</span>
								</div>
								<h6>` + user.State + `</h6>
								<a href="https://osu.ppy.sh/users/` + strconv.Itoa(user.ID) + `" target="_blank"><h2>` + user.Username + `</h2></a>
								<h4>Total PP: ` + FloatToString(user.Statistics.Pp) + `</h4>
								<h4>Global Rank: ` + strconv.Itoa(user.Statistics.Global_rank) + `</h4>
								<h4>Accuracy: ` + FloatToString(user.Statistics.Accuracy) + `</h4>
								<h4>Playcount: ` + strconv.Itoa(user.Statistics.Play_count) + `</h4>
							</div>
						</div>
					</div>
				</div>`)
	return finalString
}

func GetBackground(user player.User) string {
	if user.Background == "true" {
		return `background-image: url('` + user.CoverURL + `'); `
	}

	return ""
}
