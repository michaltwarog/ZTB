package main

import "database-tester/performance"

func main() {
	performance.RunPerformanceTest()
	performance.RunMySQLPerformanceTest()
	performance.RunOraclePerformanceTest()
}
