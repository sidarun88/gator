## Gator
Your personal blog aggregator

### Prerequisites
You will need to install the following beforehand:
- Postgresql (min v15)
- Go (at least v1.25)

### Install
Install this package using following command `go install github.com/sidarun88/gator`

### Post installation
Add a `.gatorconfig.json` file in your OS's Home folder before running gator.

### Commands
- `gator reset`: Resets gator data
- `gator register <name>`: Registers the user in gator if not found
- `gator login <name>`: Logins the registered user in gator if found
- `gator users`: Logs all the registered users to console
- `gator addfeed <url>`: Adds RSS feed to gator for current logged-in user and the logged-in user starts following the added feed
- `gator feeds`: Logs all added feeds to console
- `gator agg`: Save posts of all the feeds to gator
- `gator follow <url>`: The current logged-in user starts following the feed with url if found
- `gator unfollow <url>`: The current logged-in user stops following the feed with url if found
- `gator following`: Logs all the feeds being followed by current logged-in user
- `gator browse`: Logs all the posts to the console being followed by current logged-in user
