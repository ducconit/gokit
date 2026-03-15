package logging

type Mode string

const (
	ModeDisabled Mode = "disabled"
	ModeConsole  Mode = "console"
	ModeFile     Mode = "file"
	ModeBoth     Mode = "both"
)

type Config struct {
	Mode           Mode
	Level          string
	FilePath       string
	ConsolePretty  bool
	DisableConsole bool
}
