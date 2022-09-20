#__*__coding: utf-8__*__


import sys
import json
import os

def generate():
    
    args = sys.argv

    fn = args[1]
    path, file = os.path.split(fn)
    dstPath = "/Users/koko/go/src/BOC/press-test/result/"
    dstFile = dstPath + "/" + file.split(".")[0] + ".txt"
    f = open(fn, "r")
    output_f = open(dstFile, "w")

    lines = f.readlines()
    n_lines = len(lines)

    i = 0
    # output_f.write("Threads, Connections, Avg(latency), Max(Latency),    P50,       P75,      P90,       P99,       QPS" + "\n")
    while i < n_lines:
        line = lines[i].strip()
        i = i+1
        if line.startswith("Running"):
            new_line = []
            continue
        parts = line.split()
        if line.endswith("connections"):
            new_line.append(parts[0])
            new_line.append(parts[3])
            continue

        if line.startswith("Latency") and len(parts)>=4:
            new_line.append(parts[1])
            new_line.append(parts[3])
            continue
        if line.startswith("50%") or line.startswith("75%") or line.startswith("90%") or line.startswith("99%"):
            new_line.append(parts[1]) 
            continue

        if line.startswith("Requests/sec"):
            new_line.append(parts[1])
            s = json.dumps(new_line)
            print(s)
            output_f.write(s + "\n")
            continue

if __name__ == "__main__":
    generate()

# podAntiAffinity:
#       	  requiredDuringSchedulingIgnoredDuringExecution:
#           - podAffinityTerm:
#               labelSelector:
#                 matchExpressions:
#                 - key: app
#                   operator: In
#                   values:
#                   - yb-tserver
#                   - fp
#                   - cassandra
#               namespaces:
#               - yuga1
#               topologyKey: kubernetes.io/hostname

# podAffinity:
#           requiredDuringSchedulingIgnoredDuringExecution:
#           - labelSelector:
#               matchExpressions:
#               - key: app
#                 operator: In
#                 values:
#                 - fp
#             namespaces:
#             - yuga1
#             topologyKey: kubernetes.io/hostname 