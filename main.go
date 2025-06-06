package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func extractSubdomains(subs []string, minLabels int, filtered []string) map[string]struct{} {
	filterSet := make(map[string]struct{})
	for _, f := range filtered {
		filterSet[f] = struct{}{}
	}

	// Clean subdomains by removing filtered labels
	cleanedSubs := make([][]string, 0, len(subs))
	for _, sub := range subs {
		labels := strings.Split(sub, ".")
		var cleaned []string
		for _, label := range labels {
			if _, skip := filterSet[label]; !skip {
				cleaned = append(cleaned, label)
			}
		}
		cleanedSubs = append(cleanedSubs, cleaned)
	}

	// Count immediate subdomain hits: For each cleaned subdomain,
	// generate its parent domain by removing the leftmost label,
	// and increment count for that parent if it exists.
	domainCounts := make(map[string]int)
	for _, labels := range cleanedSubs {
		if len(labels) < 2 {
			continue // no parent domain possible
		}
		// Parent domain is labels[1:] (one level up)
		parent := strings.Join(labels[1:], ".")
		// Only count if the subdomain is strictly deeper (len(labels) > len(parent labels))
		// which is always true here as we remove one label
		domainCounts[parent]++
	}

	// Filter domains that meet minLabels in both label count and hit count
	result := make(map[string]struct{})
	for domain, count := range domainCounts {
		if count >= minLabels && len(strings.Split(domain, ".")) >= minLabels {
			result[domain] = struct{}{}
		}
	}

	return result
}

func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeOutput(subs map[string]struct{}, outputPath string) error {
	var out *os.File
	var err error

	if outputPath != "" {
		out, err = os.Create(outputPath)
		if err != nil {
			return err
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	for sub := range subs {
		fmt.Fprintln(out, sub)
	}
	return nil
}

func main() {
	inputPath := flag.String("i", "", "Input file path")
	outputPath := flag.String("o", "", "Output file path (optional)")
	minLabels := flag.Int("min", 1, "Minimum number of labels required in subdomain (default: 1)")
	filterStr := flag.String("fs", "balls", "Comma-separated subdomain labels to ignore in count (e.g. 'www,dev')")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-i input_path] [-o output_file] [-min N] [-fs labels]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *inputPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	lines, err := readLines(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}

	var filters []string
	if *filterStr != "" {
		filters = strings.Split(*filterStr, ",")
	}

	fmt.Println(*outputPath)

	subdomains := extractSubdomains(lines, *minLabels, filters)

	if err := writeOutput(subdomains, *outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}
}
