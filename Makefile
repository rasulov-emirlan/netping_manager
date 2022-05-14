# ENV stands for file which contains configs with sensitive data
ENV=.env

run:
	go run cmd/apiserver/main.go ${ENV}