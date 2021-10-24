# Medzernik Discord Bot in Go
This is a discord bot written entirely in golang using the discordgo library. 

The bot itself has a custom command parser and many custom functions. It is a featured bot you can use on your server.

# Instructions for Self-Hosting

Follow the instructions (first few steps to setup an application)
at https://www.freecodecamp.org/news/create-a-discord-bot-with-python/

Get the token and your guildID. Then invite the bot into the guild and give it permissions.

COPY THE `config.yml` FILE FROM THE EXAMPLES DIRECTORY -> ROOT NEXT TO WHERE THE MAIN EXECUTABLE IS TO USE, **THEN EDIT
IT** AND RUN THE BOT

**Currently pre-built: LINUX ARM64, LINUX AMD64, WINDOWS AMD64**

# Instructions for adding the bot directly

Coming very soon (support is already done and working!)

# Automatic functions

## Almost all functions run in separate goroutines, making the bot scaleable.

#### Planned games reminding

The bot automatically checks for any games planned every minute and reminds people the time is due

#### Trusted channel locking-unlocking

The bot automatically checks for the time to lock and unlock sendMessage permissions for a given channel and for given
roles.

# Current syntax:

**prefix:** `/` *(Discord command system)*

![image](https://user-images.githubusercontent.com/1900179/138592653-3ec239f6-a80c-481b-8960-fa513fa78acb.png)

# commands:

## muting / umuting [ADMIN] [TRUSTED]

`/mute @mention`

`/unmute @mention`

### (Muting is also available as a right-click command)
![image](https://user-images.githubusercontent.com/1900179/138592669-fd5ce32a-a7a6-4570-8816-930da9659a11.png)


mutes or unmutes a user. Checks if the user is either an Admin/Mod, or at least a Trusted user. For a trusted user, only
muting of users that have joined less than 24 hours ago is allowed.

## age check

`/age @mention`

checks the age of the user (account age).

## user join age check [ADMIN]

`/checkusers <hours to check>`

checks all the users that connected in less than 24h (by default).

## plan <hours:minutes> game_name @mentions

`/plan HH:MM game_name @MENTION`

Plans and writes a UNIX timestamped time + ID and other data into an entry in SQLite database for later use. Confirms
with outputting the entry into the Discord channel.

## planned

`/planned`

SELECT * from planned games table. Outputs into the same Discord channel.

## weather

`/weather City Name`

Check to see the weather information of a particular city. Supports cities with spaces in names. Runs concurrently.
![image](https://user-images.githubusercontent.com/1900179/138592679-b73a73c5-aeec-438f-bfea-872917751fb9.png)

## topic

`/topic`

Outputs a random topic for a discussion.

More features pending, including a docker image and a custom SQLite database for internal event planning.

## kick/ban [ADMIN]

`/kick @user <reason for the kick>`

`/ban @user <reason for the ban> <days of messages to delete>`

Kicks and or bans a user. Posts a message to the log channel defined.

In case of a ban also deletes previous 7 days of the user's messages.

DMs the user the reason and information about his kick/ban.

## version

`/version`

Displays the version of the bot

## purge [ADMIN]

`/purge NUMBERINT`

### (PurgeTo and a new command PurgeToUser available as messageRightclickCommand)
![image](https://user-images.githubusercontent.com/1900179/138592702-ea36a27f-5aa6-40d3-bc86-bb4a56875210.png)

PurgeTo rightclick deletes up to 100 messages under your selected message.

PurgeToUser does the same but filters the messages to delete only by the user of the message you clicked.

Deletes 1-100 messages in the channel that the command was typed in.

## prunecount

`/prunecount NUMBERINT`

Checks how many users would be pruned (minimum is 7 days, maximum is undefined, however, for me only 30 worked as max).

## prunemembers [ADMIN]

`/prunemembers NUMBERINT`

Prunes members that have been inactive for a set amount of days (minimum 7, max undefined, however, for me only 30
worked as max).

## members

`/members`

Counts the number of members on the server.

## configreload [ADMIN]

`/configreload`

This command reloads the config data into memory without restarting the bot.

## setuserperm [ADMIN]

`/setuserperm <allow> <@mention> <PERMID>`

Sets the channels permissions. Calculate the permissions here: https://discordapi.com/permissions.html

## setroleperm [ADMIN]

`/setchannelperm <allow> <@role> <PERMID>`

Sets the channels permissions. Calculate the permissions here: https://discordapi.com/permissions.html

## redirect [ADMIN]

`/redirect #channel`

Sets the current channel for a 360 second slowmode, and embeds a new channel for people to write to.

## slow [ADMIN]

`/slow NUMOFSECS (0-21600)`

Sets a slowmode in the current channel. 0 seconds don't do anything (bug of Discord) therefore I have set it to
autocorrect to 1.

## voicechannelmodify

`/voicechannelmodify <name> <bitrate>`

Changes the name of the channel. Optionally you can specify a bitrate. The bot tries to use the highest possible by
default.
