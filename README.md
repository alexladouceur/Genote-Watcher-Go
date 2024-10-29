# GenoteWatcher-Go

Small app that sends notifications to a discord webhook when a new note is added
or changed on genote

The original app is written in typescript but this one has been rewritten in go
to make things faster and easier to share.

## Requirements

- Create a .env file at the root with the following keys:
  - **GENOTE_USER** : Contains your UdS email to login into genote
  - **GENOTE_PASSWORD** : Contains your UdS password to login into genote
  - **DISCORD_WEBHOOK** : Your desired Discord webhook url

## Start the app

- Run the executable !

## Build the app from scratch

- You will need to have go installed
- Run the build.ps1 script. It will build a windows and a linux/amd64 executable
