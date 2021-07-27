# SlovakiaDiscordBotGo in Go
This is a discord bot written entirely in golang using the discordgo library. 

The bot itself has a custom command parser and many custom functions. It is a featured bot you can use on your server.

# Instructions for Self-Hosting
Follow the instructions (first few steps to setup an application) at https://www.freecodecamp.org/news/create-a-discord-bot-with-python/

Get the token and your guildID. Then invite the bot into the guild and give it permissions.

COPY THE `config.yml` FILE FROM THE EXAMPLES DIRECTORY -> ROOT NEXT TO WHERE THE MAIN EXECUTABLE IS TO USE, **THEN EDIT IT** AND RUN THE BOT

**Currently only LINUX ARM64 BUILDS ARE PRECOMPILED**

# Automatic functions
## These features run automatically as goroutines

#### RAID checking 
The bot automatically checks for users connected in the last ~10 minutes and alerts admins if there is a possible raid attack. This is done concurrently.
#### Planned games reminding
The bot automatically checks for any games planned every minute and reminds people the time is due
#### Trusted channel locking-unlocking
The bot automatically checks for the time to lock and unlock sendMessage permissions for a given channel and for given roles.

# Current syntax:

**prefix:** `.` *(you can change the prefix in the config)*

**command:** (no space after prefix)

**arguments:** spaces, spaced after command

**example:** `.mute @mention`

# commands:
## muting / umuting [ADMIN] [TRUSTED]
`.mute @mention @mention @mention ...`

`.unmute @mention @mention @mention ...`

mutes or unmutes a user. Checks if the user is either an Admin/Mod, or at least a Trusted user. For a trusted user, only muting of users that have joined less than 24 hours ago is allowed.

## age check
`.age @mention`

checks the age of the user (account age).

## user join age check [ADMIN]
`.checkusers`

checks all the users that connected in less than 24h.

## plan <hours:minutes> game_name @mentions
`.plan HH:MM game_name @MENTION`

Plans and writes a UNIX timestamped time + ID and other data into an entry in SQLite database for later use. Confirms with outputting the entry into the Discord channel.

## planned
`.planned`

SELECT * from planned games table. Outputs into the same Discord channel.

## weather
`.weather City Name String`

Check to see the weather information of a particular city. Supports cities with spaces in names. Runs concurrently.

## topic
`.topic`

Outputs a random topic for a discussion.

More features pending, including a docker image and a custom SQLite database for internal event planning.

## kick/ban [ADMIN]
`.kick @user <reason for the kick>`

`.ban @user <reason for the ban>`

Kicks and or bans a user. Posts a message to the log channel defined. 

In case of a ban also deletes previous 7 days of the user's messages.

DMs the user the reason and information about his kick/ban.

## version
`.version`

Displays the version of the bot

## purge [ADMIN]
`.purge NUMBERINT`

Deletes 1-100 messages in the channel that the command was typed in.

## prunecount
`.prunecount NUMBERINT`

Checks how many users would be pruned (minimum is 7 days, maximum is undefined, however, for me only 30 worked as max).

## prunemembers [ADMIN]
`.prunemembers NUMBERINT`

Prunes members that have been inactive for a set amount of days (minimum 7, max undefined, however, for me only 30 worked as max).

## members
`.members`

Counts the number of members on the server.

## configreload [ADMIN]
`.configreload`

This command reloads the config data into memory without restarting the bot.