import sys

with open('cmd/helios/main.go', 'r') as f:
    lines = f.readlines()

new_lines = []
skip = 0
for i, line in enumerate(lines):
    if skip > 0:
        skip -= 1
        continue
    
    # Lines 16-19 are:
    # 16: func main() {
    # 17: )
    # 18: (blank)
    # 19: func main() {
    if i == 15: # Line 16
        new_lines.append("func main() {\n")
        skip = 3
        continue
    
    # 25: switch os.Args[1] {
    # 26: case "hash":
    # 27: case "--version", "-v":
    # 28:         fmt.Printf("helios %s\n", version)
    # 29:         return
    # 30: case "hash":
    if i == 25: # Line 26: case "hash":
        continue # skip first duplicate hash case
        
    new_lines.append(line)

# Update printUsage
final_lines = []
for line in new_lines:
    final_lines.append(line)
    if "helios verify <vectors.json>" in line:
        final_lines.append('        fmt.Fprintln(os.Stderr, "  helios --version             Show version")\n')

with open('cmd/helios/main.go', 'w') as f:
    f.writelines(final_lines)
