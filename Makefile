test_coverage:
	PWD=pwd
	export GOPATH=$(PWD):$(GOPATH)
	echo $(GOPATH)
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
