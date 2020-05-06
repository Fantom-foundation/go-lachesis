.PHONY: mock-backend
mock-backend:
	rm -f ./backend_mock_test.go
	mockgen -package=ethapi -self_package=github.com/Fantom-foundation/go-lachesis/ethapi -source=backend.go github.com/Fantom-foundation/go-lachesis/ethapi Backend > ./_backend_mock_test.go
	mv -f ./_backend_mock_test.go ./backend_mock_test.go
