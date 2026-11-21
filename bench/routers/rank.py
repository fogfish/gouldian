import re,sys,collections,statistics as st
d=collections.defaultdict(lambda: collections.defaultdict(list))
for l in open(sys.argv[1]):
    m=re.match(r'Benchmark(\w+?)_(\w+?)(?:-\d+)?\s+\d+\s+([\d.]+) ns/op\s+(\d+) B/op\s+(\d+) allocs',l)
    if m: d[m.group(2)][m.group(1)].append((float(m.group(3)),int(m.group(4)),int(m.group(5))))
cats=sys.argv[2:] or sorted(d)
for c in cats:
    rows=sorted(((st.median(x[0] for x in v),k,v[0][1],v[0][2]) for k,v in d[c].items()))
    print(f"\n{c}")
    for i,(ns,k,b,a) in enumerate(rows,1):
        print(f"  {i:2d}. {k:12s} {ns:10.1f} ns  {b:6d} B  {a:4d} allocs")
