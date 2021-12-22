build:
	go build
	mkdir fixtures || exit 0
	mkdir ./fixtures/tracker || exit 0
	mv ./tracker ./fixtures/tracker/

cli-test: build
	cucumber

mocks:
	mockgen -source=domain/inbox_repo.go -destination=domain/mock_inbox_repo.go -package=domain
	mockgen -source=domain/session_repo.go -destination=domain/mock_session_repo.go -package=domain