package loggingdrain

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData = readTestData()

func BenchmarkBuildTree(b *testing.B) {
	miner, _ := NewTemplateMiner()
	for i := 0; i < b.N; i++ {
		miner.AddLogMessage(testData[i%len(testData)])
	}
}

func BenchmarkBuildDepth5Tree(b *testing.B) {
	miner, _ := NewTemplateMiner(WithDrainDepth(5))
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

func BenchmarkMatchDepth5Tree(b *testing.B) {
	miner, _ := NewTemplateMiner(WithDrainDepth(5))
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

func TestToJson(t *testing.T) {
	t.Run("test to json", func(t *testing.T) {
		testJson := `{"Drain":{"MaxDepth":4,"Sim":0.4,"MaxChildren":100,"MaxClusters":1000,"ClusterCounter":2,"Clusters":[{"ID":1,"LogTemplateTokens":["Dec","10","[*]","LabSZ","[*]","input_userauth_request:","invalid","user","[*]","[preauth]"]},{"ID":2,"LogTemplateTokens":["Dec","10","[*]","LabSZ","[*]","Failed","password","for","invalid","user","[*]","from","0.0.0.0","port","[*]","ssh2"]}],"RootNode":{"NodeType":0,"Length":0,"TokenNodeChildren":{},"LengthNodeChildren": {"10":{"NodeType":1,"Length":10,"TokenNodeChildren":{"Dec":{"NodeType":2,"Length":0,"TokenNodeChildren":{},"LengthNodeChildren":{},"Clusters":[{"ID":1, "LogTemplateTokens":["Dec","10","[*]","LabSZ","[*]","input_userauth_request:","invalid","user","[*]","[preauth]"]}]}},"LengthNodeChildren":{},"Clusters":[]},"16":{"NodeType":1,"Length":16,"TokenNodeChildren":{"Dec":{"NodeType":2,"Length":0,"TokenNodeChildren":{},"LengthNodeChildren":{},"Clusters":[{"ID":2,"LogTemplateTokens":["Dec","10","[*]","LabSZ","[*]","Failed","password","for","invalid","user","[*]","from","0.0.0.0","port","[*]","ssh2"]}]}},"LengthNodeChildren":{},"Clusters":[]}},"Clusters":[]}},"Masker":{"Prefix":"[:","Suffix":":]","MaskInstructions":[{"Pattern":"abc","MaskWith":"abc"}]}}`

		miner, _ := NewTemplateMiner(WithMaskInsturction("abc", "abc"))
		rawLogs := []string{
			"Dec 10 07:07:38 LabSZ sshd[24206]: input_userauth_request: invalid user test9 [preauth]",
			"Dec 10 07:08:28 LabSZ sshd[24208]: input_userauth_request: invalid user webmaster [preauth]",
			"Dec 10 09:12:32 LabSZ sshd[24490]: Failed password for invalid user ftpuser from 0.0.0.0 port 62891 ssh2",
			"Dec 10 09:12:35 LabSZ sshd[24492]: Failed password for invalid user pi from 0.0.0.0 port 49289 ssh2",
			"Dec 10 09:12:44 LabSZ sshd[24501]: Failed password for invalid user ftpuser from 0.0.0.0 port 60836 ssh2",
			"Dec 10 07:28:03 LabSZ sshd[24245]: input_userauth_request: invalid user pgadmin [preauth]",
		}
		for _, log := range rawLogs {
			miner.AddLogMessage(log)
		}
		b, err := json.Marshal(miner)
		if err != nil {
			t.Fatal(err)
		}
		assert.JSONEq(t, testJson, string(b), "json should match")

		newMiner := TemplateMiner{}
		if err := json.Unmarshal(b, &newMiner); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, miner, &newMiner, "miner should match")
	})
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
