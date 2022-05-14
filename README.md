# netping_manager

This is a manager which provides a simple web ui to control netpings. It is a small project i built for an internship.

## How to use it?

- Migrate all the tables from /migrations if needed
- Compile it with `go build cmd/apiserver/main.go`
- Run the binary with `<binary name> <name of the config file>`. The config file needs to be a .env file. If you do not provide the name of the file, then server will try to get config data from enviorment variables. Example of .env file can be seen in the .`/example.dev` file

## Important

- This application is not just an api. It also serves a website. That website was built with `react` and its files should be stored in `./cmd/apiserver/dist/` directory. If you will modify the client, put the newly generated website in that directory. And remember that our app can only serve SPA websites. And base path for them should be `/website`.
- In this app only admins are allowed to add new users. So if you are deploying it for the first time, just add your first user manualy. You can use `./password_gen.go` for generating a hashed password for it. And do not forget to make that user and admin.
