# netping_manager
This is a manager which provides a simple web ui to control netpings. It is a small project i built for an internship.

## How to use it?
- Migrate all the tables from /migrations
- Compile it with ```go build cmd/apiserver/main.go```
- Run the binary with ```<binary name> <name of the config file>```. The config file needs to be a .env file. If you do not provide the name of the file, then server will try to get config data from enviorment variables.
