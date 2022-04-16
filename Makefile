DB_URL='mysql://mysql:root@tcp(localhost:3306)/netping_manager'

migrate_create:
	migrate create -ext sql -dir ./migrations -seq ${TABLE_NAME}

migrate_up:
	migrate -database ${DB_URL} -path ./migrations/ up 

migrate_down:
	migrate -database ${DB_URL} -path ./migrations/ down

migrate_force_fix:
	migrate -path ./migrations/ -database ${DB_URL} force ${VERSION}

mock_repository:
	mockgen -destination=internal/app/services/mocks/${NAME}.go \
	-package=mocks 	-source=internal/app/services/${NAME}.go \
	github.com/sea-auca/auca-issue-collector/internal/app/services/${NAME}.go ${NAME}Repository


setup:
	go mod tidy
	go install github.com/golang/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
