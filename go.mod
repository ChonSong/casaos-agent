module github.com/ChonSong/casaos-agent

go 1.22.0

require (
	github.com/ChonSong/casaos-agent/CasaOS-CLI v0.0.0
	github.com/spf13/cobra v1.9.1
	github.com/spf13/viper v1.19.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/ChonSong/casaos-agent/CasaOS-CLI => ./CasaOS-CLI
