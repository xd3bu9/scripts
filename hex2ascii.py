#!/usr/bin/env python

import sys

hexval = ''

for line in sys.stdin:
	hexval += line

print(chr(int(hexval.rstrip('\n\r'), 16)))

# script to convert hexadecimal values to their ascii representation.
# usage: echo 0x41 | hex2ascii