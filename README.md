# deeperDomainExtractor
Extract deeper level domains of a target from a list of subdomains

# Usage
```
Usage: ./deeperDomainExtractor [-i input_path] [-o output_file] [-min N] [-fs labels]
  -fs string
        Comma-separated subdomain labels to ignore in count (e.g. 'www,dev')
  -i string
        Input file path
  -min int
        Minimum number of labels required in subdomain (default: 1)
  -o string
        Output file path
```

# Example

using a wordlist of 

```
www.test.example
www.hello.example
www.b.hello.example
aaa.b.hello.example
ccc.b.hello.example
what.hello.example
```

```
./deeperDomainExtractor -fs www,dev -i test.txt -min 2

hello.example
b.hello.example

./deeperDomainExtractor -fs dev -i test.txt -min 1

test.example
hello.example
b.hello.example
```