BIN := `pwd | awk -F"/" '{print $$NF}'`
SRC := `if [ -d entry/ ]; then echo entry; else echo entry; fi`

default: release

gensrc:
		./version $(SRC)/version.go

release: gensrc
		go build -i -o ./bin/$(BIN) ./$(SRC)
race:
		go build -i -v -o ./bin/$(BIN) -race ./$(SRC)
run:
		./bin/$(BIN)
stop:
		./bin/$(BIN) -s stop
w:
		go get github.com/codeskyblue/fswatch && fswatch
dev:
		go build -gcflags "-N -l" -i -x -o ./bin/$(BIN) ./$(SRC) && killall goserver && ./bin/$(BIN) -c ../conf/conf_dev.ini

