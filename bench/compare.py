#!/usr/bin/env python3
"""Compare bench/v2.txt with bench/v3.txt (median ns/op): python3 bench/compare.py"""
import re, collections, statistics as st

def load(p):
    d = collections.defaultdict(list)
    for l in open(p):
        m = re.match(r'(Benchmark\S+)\s+\d+\s+([\d.]+) ns/op\s+(\d+) B/op\s+(\d+) allocs', l)
        if m:
            d[m.group(1)].append((float(m.group(2)), int(m.group(3)), int(m.group(4))))
    return d

a, b = load('bench/v2.txt'), load('bench/v3.txt')
print("%-40s %9s %9s %7s  allocs" % ("benchmark", "v2 ns", "v3 ns", "ratio"))
for k in sorted(a):
    if k in b:
        x = st.median(v[0] for v in a[k]); y = st.median(v[0] for v in b[k])
        print("%-40s %9.1f %9.1f %6.2fx  %d -> %d" % (k, x, y, y / x, a[k][0][2], b[k][0][2]))
