package command

type CommandStatus struct {
	OK      string
	ERR     string
	SYNTAX  string
	WARN    string
	AUTH    string
	AUTOFIX string
}

// StatusBot StatusBot is a variable to pass to the messageEmbed to make an emoji
var StatusBot CommandStatus = CommandStatus{
	OK:      "",
	ERR:     ":bangbang: ERROR",
	SYNTAX:  ":question: SYNTAX",
	WARN:    ":warning: WARNING",
	AUTH:    ":no_entry: AUTHENTICATION",
	AUTOFIX: ":wrench: AUTOCORRECTING",
}
