.PHONY: all clean test

all: test

clean:
	rm -rf */example.git */checkout */terraform.tfstate.backup */terraform.tfstate */plan.out */destroy.out

test: clean
	for i in $$(find . -maxdepth 2 -name '*_test.go' -type f); do cd $$(dirname $$i); go test; cd ..; done
	