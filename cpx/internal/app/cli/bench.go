package cli

import (
	"fmt"

	"github.com/ozacod/cpx/internal/pkg/naming"
)

type benchSources struct {
	Main string
}

func generateBenchmarkArtifacts(projectName string, bench string) (*benchSources, []string) {
	switch bench {
	case "google-benchmark":
		return &benchSources{Main: googleBenchMain(projectName)}, []string{"benchmark"}
	case "nanobench":
		return &benchSources{Main: nanoBenchMain(projectName)}, []string{"nanobench"}
	case "catch2-benchmark":
		return &benchSources{Main: catch2BenchMain(projectName)}, []string{"catch2"}
	default:
		return nil, nil
	}
}

func googleBenchMain(projectName string) string {
	safeName := naming.SafeIdent(projectName)
	return fmt.Sprintf(`#include <benchmark/benchmark.h>
#include <%s/%s.hpp>

static void BM_version(benchmark::State& state) {
    for (auto _ : state) {
        benchmark::DoNotOptimize(%s::version());
    }
}

BENCHMARK(BM_version);

int main(int argc, char** argv) {
    benchmark::Initialize(&argc, argv);
    if (benchmark::ReportUnrecognizedArguments(argc, argv)) return 1;
    benchmark::RunSpecifiedBenchmarks();
}
`, projectName, projectName, safeName)
}

func nanoBenchMain(projectName string) string {
	safeName := naming.SafeIdent(projectName)
	return fmt.Sprintf(`#include <nanobench.h>
#include <%s/%s.hpp>
#include <iostream>

int main() {
    ankerl::nanobench::Bench bench;
    bench.run("version", [] {
        ankerl::nanobench::doNotOptimizeAway(%s::version());
    });
    return 0;
}
`, projectName, projectName, safeName)
}

func catch2BenchMain(projectName string) string {
	safeName := naming.SafeIdent(projectName)
	return fmt.Sprintf(`#include <catch2/catch_all.hpp>
#include <%s/%s.hpp>

TEST_CASE("Benchmark version", "[benchmark]") {
    BENCHMARK("version") {
        return %s::version();
    };
}
`, projectName, projectName, safeName)
}
