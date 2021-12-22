build:
	go build

mocks:
	mockgen -source=domain/inbox_repo.go -destination=domain/mock_inbox_repo.go -package=domain
	mockgen -source=domain/session_repo.go -destination=domain/mock_session_repo.go -package=domain