test:
	go test github.com/golang_samples/golang_test/number --cover -v
test-coverage:
	go test github.com/golang_samples/golang_test/number --cover -coverprofile=coverage.out
	go tool cover -html=coverage.out
# Patter test case match "Bench" only
bechmark:
	go test github.com/golang_samples/golang-benchmark -run=Bench -bench=.
# Patter test case match "Calc" only
bechmark-calc:
	go test github.com/golang_samples/golang-benchmark -run=Calc -bench=.