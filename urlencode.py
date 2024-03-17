#!/usr/bin/env python

import sys
import urllib.parse

unencoded = ''

for line in sys.stdin:
	unencoded += line

print(urllib.parse.quote(unencoded.rstrip('\n\r'), safe=''), end='')

# python script to urlencode data from stdin
# usage: echo "@+data" | urlencode
# credit: 0xtib3rius