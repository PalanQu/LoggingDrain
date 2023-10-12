package loggingdrain

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
)

var testData = readTestData()

func BenchmarkBuildTree(b *testing.B) {
	miner, _ := NewTemplateMiner()
	for i := 0; i < b.N; i++ {
		miner.AddLogMessage(testData[i%len(testData)])
	}
}

func BenchmarkMatchTree(b *testing.B) {
	miner, _ := NewTemplateMiner()
	for _, log := range testData {
		miner.AddLogMessage(log)
	}
	for i := 0; i < b.N; i++ {
		miner.Match(testData[i%len(testData)])
	}
}

func BenchmarkUnmarshalJson(b *testing.B) {
	testJson := `{"timestamp":"Jun 22 04:30:55","host":"combo","process":"sshd","pid":"pam_unix","log_level":"17125","message":"authentication failure","uid":"0","euid":"0","tty":"NODEVssh","ruser":"","rhost":"ip-216-69-169-168.ip.secureserver.net"}`
	data := map[string]string{}
	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(testJson), &data)
	}
}

func readTestData() []string {
	testData := []string{}
	logFile, err := os.Open("test_data/Linux_2k.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		line := scanner.Text()
		testData = append(testData, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return testData
}
