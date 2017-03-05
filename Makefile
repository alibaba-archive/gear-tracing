test:
	go test --race

cover:
	rm -f *.coverprofile
	go test -coverprofile=tracing.coverprofile
	go tool cover -html=tracing.coverprofile
	rm -f *.coverprofile

.PHONY: test cover
