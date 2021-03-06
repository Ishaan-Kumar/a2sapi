package config

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/syncore/a2sapi/src/constants"
	"github.com/syncore/a2sapi/src/steam/filters"
)

const (
	defaultMaxHostsToReceive        = 4000
	defaultAutoQueryMaster          = true
	defaultTimeBetweenMasterQueries = 90
	defaultUseWebServerList         = true
	// defaultTimeForHighServerCount: not used in JSON, only in the config dialog
	defaultTimeForHighServerCount = 120
)

// CfgSteam represents Steam-related configuration options.
type CfgSteam struct {
	AutoQueryMaster          bool   `json:"timedMasterServerQuery"`
	SteamWebAPIKey           string `json:"steamWebAPIKey"`
	UseWebServerList         bool   `json:"useWebServerList"`
	AutoQueryGame            string `json:"gameForTimedMasterQuery"`
	TimeBetweenMasterQueries int    `json:"timeBetweenMasterQueries"`
	MaximumHostsToReceive    int    `json:"maxHostsToReceive"`
}

func configureTimedMasterQuery(reader *bufio.Reader) bool {
	valid, val := false, false
	prompt := fmt.Sprintf(`
Perform an automatic retrieval of game servers from the Steam master server at
timed intervals? This is necessary if you want the API to maintain a filterable / searchable
list of game servers.
%s`, promptColor("> 'yes' or 'no' [default: %s]: ",
		getBoolString(defaultAutoQueryMaster)))

	input := func(r *bufio.Reader) (bool, error) {
		enable, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultAutoQueryMaster,
				fmt.Errorf("Unable to read respone: %s", rserr)
		}
		if enable == newline {
			return defaultAutoQueryMaster, nil
		}
		response := strings.Trim(enable, newline)
		if strings.EqualFold(response, "y") || strings.EqualFold(response, "yes") {
			return true, nil
		} else if strings.EqualFold(response, "n") || strings.EqualFold(response,
			"no") {
			return false, nil
		} else {
			return defaultAutoQueryMaster,
				fmt.Errorf("[ERROR] Invalid response. Valid responses: y, yes, n, no")
		}
	}
	var err error
	for !valid {
		fmt.Fprintf(color.Output, prompt)
		val, err = input(reader)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureUseSteamWebAPIList(reader *bufio.Reader) bool {
	valid, val := false, false
	prompt := fmt.Sprintf(`
Use the game server list provided by the Steam Web API? This can be considerably faster than having
the application perform queries to the master server, since it will use the most up-to-date
server list that has already been processed by Valve. This is also preferable because the master
server at times has random reliability and downtime issues on Valve's end. Note, without this method
Valve will throttle your requests if more than 30 UDP packets sent (game has 6930 or more servers)
within 60 seconds.
%s`, promptColor("> 'yes' or 'no' [default: %s]: ",
		getBoolString(defaultUseWebServerList)))

	input := func(r *bufio.Reader) (bool, error) {
		enable, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultUseWebServerList,
				fmt.Errorf("Unable to read respone: %s", rserr)
		}
		if enable == newline {
			return defaultUseWebServerList, nil
		}
		response := strings.Trim(enable, newline)
		if strings.EqualFold(response, "y") || strings.EqualFold(response, "yes") {
			return true, nil
		} else if strings.EqualFold(response, "n") || strings.EqualFold(response,
			"no") {
			return false, nil
		} else {
			return defaultUseWebServerList,
				fmt.Errorf("[ERROR] Invalid response. Valid responses: y, yes, n, no")
		}
	}
	var err error
	for !valid {
		fmt.Fprintf(color.Output, prompt)
		val, err = input(reader)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureSteamWebAPIKey(reader *bufio.Reader) string {
	valid := false
	var val string
	prompt := fmt.Sprintf(`
Enter your Steam Web API key for querying Valve's server list. You can receive a Web
API key for free from http://steamcommunity.com/dev/apikey
%s`, promptColor("> [default: NONE]: "))

	input := func(r *bufio.Reader) (string, error) {
		keyval, rserr := r.ReadString('\n')
		if rserr != nil {
			return "", fmt.Errorf("Unable to read respone: %s", rserr)
		}
		if keyval == newline {
			return "", fmt.Errorf("[ERROR] Invalid response. Valid response is an API key")
		}
		response := strings.Trim(keyval, newline)
		return response, nil
	}
	var err error
	for !valid {
		fmt.Fprintf(color.Output, prompt)
		val, err = input(reader)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureTimedQueryGame(reader *bufio.Reader) string {
	valid := false
	var val string
	games := strings.Join(filters.GetGameNames(), "\n")
	prompt := fmt.Sprintf(`
Choose the game you would like to automatically retrieve servers for at timed
intervals. Possible choices are:
%s
More games can be added via the %s file.
%s`, games, constants.GameFileFullPath, promptColor("> [default: NONE]: "))

	input := func(r *bufio.Reader) (string, error) {
		gameval, rserr := r.ReadString('\n')
		if rserr != nil {
			return "", fmt.Errorf("Unable to read respone: %s", rserr)
		}
		if gameval == newline {
			return "", fmt.Errorf("[ERROR] Invalid response. Valid responses:\n%s", games)
		}
		response := strings.Trim(gameval, newline)
		if filters.IsValidGame(response) {
			// format the capitalization
			return filters.GetGameByName(response).Name, nil
		}
		return "", fmt.Errorf("[ERROR] Invalid response. Valid responses: %s", games)
	}
	var err error
	for !valid {
		fmt.Fprintf(color.Output, prompt)
		val, err = input(reader)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureMaxServersToRetrieve(reader *bufio.Reader) int {
	valid := false
	var val int
	prompt := fmt.Sprintf(`
Enter the maximum number of servers to retrieve from the Steam Master Server at
a time. This can be no more than 6930.
%s`, promptColor("> [default: %d]: ", defaultMaxHostsToReceive))

	input := func(r *bufio.Reader) (int, error) {
		hostsval, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultMaxHostsToReceive, fmt.Errorf("Unable to read response: %s", rserr)
		}
		if hostsval == newline {
			return defaultMaxHostsToReceive, nil
		}
		response, rserr := strconv.Atoi(strings.Trim(hostsval, newline))
		if rserr != nil {
			return defaultMaxHostsToReceive,
				fmt.Errorf("[ERROR] Maximum hosts to receive from master server must be between 500 and 6930")
		}
		if response < 500 || response > 6930 {
			return defaultMaxHostsToReceive,
				fmt.Errorf("[ERROR] Maximum hosts to receive from master server must be between 500 and 6930")
		}
		return response, nil
	}
	var err error
	for !valid {
		fmt.Fprintf(color.Output, prompt)
		val, err = input(reader)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureTimeBetweenQueries(reader *bufio.Reader, game string) int {
	valid := false
	var val int
	defaultVal := defaultTimeBetweenMasterQueries
	highServerCountGame := filters.HasHighServerCount(game)
	if highServerCountGame {
		defaultVal = defaultTimeForHighServerCount
	}
	retrievalMethodMsg := `Note: if the game returns more than 6930 servers/min,
Valve will throttle future requests for 1 min.`
	if defaultUseWebServerList {
		retrievalMethodMsg = ""
	}
	prompt := fmt.Sprintf(`
Enter the time, in seconds, between requests to grab all servers from the master
server. For many games this needs to be at least 60. For some games this will
need to be even higher. %s
%s `, retrievalMethodMsg, promptColor("> [default: %d]: ", defaultVal))

	input := func(r *bufio.Reader) (int, error) {
		timeval, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultVal,
				fmt.Errorf("Unable to read response: %s", rserr)
		}
		if timeval == newline {
			return defaultVal, nil
		}
		response, rserr := strconv.Atoi(strings.Trim(timeval, newline))
		if rserr != nil {
			return defaultVal,
				fmt.Errorf("[ERROR] Time between Steam aster server queries must be at least 60")
		}
		if response < 60 {
			return defaultVal,
				fmt.Errorf("[ERROR] Time between Steam master server queries must be at least 60")
		}
		if highServerCountGame && response < defaultTimeForHighServerCount {
			return defaultVal, fmt.Errorf(`
[ERROR] Game %s typically returns more than 6930 servers so the time between
Steam master server queries will need to be at least %d`, game,
				defaultTimeForHighServerCount)
		}
		return response, nil
	}
	var err error
	for !valid {
		fmt.Fprintf(color.Output, prompt)
		val, err = input(reader)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}
