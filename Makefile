fmt:
	for d in . internal/*; do go fmt $$d/*.go; done

install:
	go install litebrew.go
