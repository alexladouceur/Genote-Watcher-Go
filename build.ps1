# Build for Linux
$Env:GOOS = "linux"; $Env:GOARCH = "amd64"
go build -ldflags "-w -s"

# Build for Windows
$Env:GOOS = "windows"; $Env:GOARCH = "amd64"
go build -ldflags "-w -s"
