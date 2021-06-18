# SlovakiaDiscordBotGo in Go



This is a discord bot written entirely in golang using the discordgo library. The bot itself has a custom command parser and many custom functions. It serves as an example of a bot that can do very basic commands.

#Automatic functions:
#### RAID checking 
The bot automatically checks for users connected in the last ~10 minutes and alerts admins if there is a possible raid attack. This is done concurrently.
#### Planned games reminding
The bot automatically checks for any games planned every minute and reminds people the time is due

# Current syntax:

**prefix:** .

**command:** ni space after prefix

**arguments:** spaces, spaced after command

**example:** .mute @user

# commands:
## muting
.mute @user

mutes a user, except this bot and except the priviledged admin role ID user. admin power checking TBD.

## age check
.age @user

checks the age of the user (account age).

## user join age check
.check-users

checks all the users that connected in less than 24h.

## plan <hours:minutes> gamename @mentions
.plan 10:40 terraria @medzernik

Plans and writes a UNIX timestamped time + ID and other data into an entry in SQLite database for later use. Confirms with outputting the entry into the Discord channel.

## planned
.planned

SELECT * from planned games table. Outputs into the same Discord channel.



More features pending, including a docker image and a custom SQLite database for internal event planning.
