# twitch window

this is go code monitors a twitch chat

## setup

three environment variables:

TWITCH_USERNAME
: This is the user we're connecting as.

TWITCH_TOKEN
: This is the twitch oauth token.   You can obtain it from this site:
: [twitchapps.com/tmi/](https://twitchapps.com/tmi/) , the value should look like this:
: "oauth:7dplkxIP63NGrjEGunKdS3hcpkAige

- TWITCH_CHANNEL (optional)
: Channel to join, will default to TWITCH_USERNAME if not set.


