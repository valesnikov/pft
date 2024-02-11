PREFIX=/usr/local
DESTDIR=
SOURCES=pft.go
EXECUTABLE=pft

TEST_DIR=test

all: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	go build -o $(EXECUTABLE) $(SOURCES)

install: $(EXECUTABLE)
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp $(EXECUTABLE) $(DESTDIR)$(PREFIX)/bin

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/$(EXECUTABLE)

clean:
	rm $(EXECUTABLE)
	rm -rf $(TEST_DIR)

test: $(EXECUTABLE)
	mkdir -p $(TEST_DIR)/in
	mkdir -p $(TEST_DIR)/out

	dd if=/dev/random of=$(TEST_DIR)/in/big.rnd bs=1M count=512 status=progress
	dd if=/dev/random of=$(TEST_DIR)/in/small2.rnd bs=1K count=1
	dd if=/dev/random of=$(TEST_DIR)/in/small3.rnd bs=1K count=2
	dd if=/dev/random of=$(TEST_DIR)/in/small4.rnd bs=1K count=3
	dd if=/dev/random of=$(TEST_DIR)/in/small5.rnd bs=1K count=4

	./pft hs 12338 $(TEST_DIR)/in/* &
	./pft cr localhost 12338 $(TEST_DIR)/out/
	sleep 1

	#rm $(TEST_DIR)/out/small3.rnd

	diff $(TEST_DIR)/in/ $(TEST_DIR)/out/ || (echo "hs and cr failed $$?"; exit 1)
	rm $(TEST_DIR)/out/*

	./pft hr 12339 $(TEST_DIR)/out/ &
	./pft cs localhost 12339 $(TEST_DIR)/in/*
	sleep 1

	diff $(TEST_DIR)/in/ $(TEST_DIR)/out/ || (echo "hr and cs failed $$?"; exit 1)
	rm $(TEST_DIR)/out/*

	rm -rf $(TEST_DIR)
	echo "Done"
