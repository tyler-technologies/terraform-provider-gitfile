.PHONY: all clean test

all: test

clean:
	rm -rf */example.git */checkout */terraform.tfstate.backup */terraform.tfstate */plan.out */destroy.out

test: clean
	go test ./...
	