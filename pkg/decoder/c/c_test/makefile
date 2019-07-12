cc=gcc
CFLAGS=-I.. -L..
LIBS=-lasncodec

all: test_encode test_decode

test_encode: test_encode.o
	$(CC) $(CFLAGS) $(LIBS) -o $@ $<

test_decode: test_decode.o
	$(CC) $(CFLAGS) $(LIBS) -o $@ $<

.PHONY: clean

clean:
	rm *.o test_encode test_decode