package templateminer

import (
	"fmt"
	"loggingdrain/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddLogMessage(t *testing.T) {
	t.Run("test add message", func(t *testing.T) {
		rawLogs := []string{
			"Dec 10 07:07:38 LabSZ sshd[24206]: input_userauth_request: invalid user test9 [preauth]",
			"Dec 10 07:08:28 LabSZ sshd[24208]: input_userauth_request: invalid user webmaster [preauth]",
			"Dec 10 09:12:32 LabSZ sshd[24490]: Failed password for invalid user ftpuser from 0.0.0.0 port 62891 ssh2",
			"Dec 10 09:12:35 LabSZ sshd[24492]: Failed password for invalid user pi from 0.0.0.0 port 49289 ssh2",
			"Dec 10 09:12:44 LabSZ sshd[24501]: Failed password for invalid user ftpuser from 0.0.0.0 port 60836 ssh2",
			"Dec 10 07:28:03 LabSZ sshd[24245]: input_userauth_request: invalid user pgadmin [preauth]",
		}
		expected := []string{
			"Dec 10 07:07:38 LabSZ sshd[24206]: input_userauth_request: invalid user test9 [preauth]",
			"Dec 10 [*] LabSZ [*] input_userauth_request: invalid user [*] [preauth]",
			"Dec 10 09:12:32 LabSZ sshd[24490]: Failed password for invalid user ftpuser from 0.0.0.0 port 62891 ssh2",
			"Dec 10 [*] LabSZ [*] Failed password for invalid user [*] from 0.0.0.0 port [*] ssh2",
			"Dec 10 [*] LabSZ [*] Failed password for invalid user [*] from 0.0.0.0 port [*] ssh2",
			"Dec 10 [*] LabSZ [*] input_userauth_request: invalid user [*] [preauth]",
		}
		expectedClusterSize := []int{
			1, 2, 1, 2, 3, 3,
		}
		drain := newDrain()
		for i, rawLog := range rawLogs {
			cluster, _ := drain.addLogMessage(rawLog)
			assert.Equal(t, expected[i], cluster.getTemplate())
			assert.Equal(t, expectedClusterSize[i], len(cluster.logs))
		}
		// fmt.Println(drain.status())
	})
	t.Run("test empty msg", func(t *testing.T) {
		drain := newDrain()
		for i := 0; i < 3; i++ {
			drain.addLogMessage("")
		}
		assert.Equal(t, 1, len(drain.idToCluster.Keys()))
	})
	t.Run("add short message", func(t *testing.T) {
		model := newDrain()
		_, updateType := model.addLogMessage("hello")
		assert.Equal(t, CLUSTER_UPDATE_TYPE_NEW_CLUSTER, updateType)
		_, updateType = model.addLogMessage("hello")
		assert.Equal(t, CLUSTER_UPDATE_TYPE_NONE, updateType)
		_, updateType = model.addLogMessage("otherword")
		assert.Equal(t, CLUSTER_UPDATE_TYPE_NEW_CLUSTER, updateType)
	})
	t.Run("add log message sim75", func(t *testing.T) {
		rawLogs := []string{
			"Dec 10 07:07:38 LabSZ sshd[24206]: input_userauth_request: invalid user test9 [preauth]",
			"Dec 10 07:08:28 LabSZ sshd[24208]: input_userauth_request: invalid user webmaster [preauth]",
			"Dec 10 09:12:32 LabSZ sshd[24490]: Failed password for invalid user ftpuser from 0.0.0.0 port 62891 ssh2",
			"Dec 10 09:12:35 LabSZ sshd[24492]: Failed password for invalid user pi from 0.0.0.0 port 49289 ssh2",
			"Dec 10 09:12:44 LabSZ sshd[24501]: Failed password for invalid user ftpuser from 0.0.0.0 port 60836 ssh2",
			"Dec 10 07:28:03 LabSZ sshd[24245]: input_userauth_request: invalid user pgadmin [preauth]",
		}
		expected := []string{
			"Dec 10 07:07:38 LabSZ sshd[24206]: input_userauth_request: invalid user test9 [preauth]",
			"Dec 10 07:08:28 LabSZ sshd[24208]: input_userauth_request: invalid user webmaster [preauth]",
			"Dec 10 09:12:32 LabSZ sshd[24490]: Failed password for invalid user ftpuser from 0.0.0.0 port 62891 ssh2",
			"Dec 10 [*] LabSZ [*] Failed password for invalid user [*] from 0.0.0.0 port [*] ssh2",
			"Dec 10 [*] LabSZ [*] Failed password for invalid user [*] from 0.0.0.0 port [*] ssh2",
			"Dec 10 07:28:03 LabSZ sshd[24245]: input_userauth_request: invalid user pgadmin [preauth]",
		}
		drain := newDrain(withSim(0.75))
		for i, rawLog := range rawLogs {
			cluster, _ := drain.addLogMessage(rawLog)
			assert.Equal(t, expected[i], cluster.getTemplate())
		}
	})
	t.Run("test maxCluster", func(t *testing.T) {
		rawLogs := []string{
			"A format 1",
			"A format 2",
			"B format 1",
			"B format 2",
			"A format 3",
		}
		expected := []string{
			"A format 1",
			"A format [*]",
			"B format 1",
			"B format [*]",
			"A format 3",
		}
		drain := newDrain(withMaxClusters(1))
		for i, rawLog := range rawLogs {
			cluster, _ := drain.addLogMessage(rawLog)
			assert.Equal(t, expected[i], cluster.getTemplate())
		}
	})
	t.Run("test lru", func(t *testing.T) {
		rawLogs := []string{
			"A A A",
			"A A B",
			"A B A",
			"A B B",
			"A C A",
			"A C B",
			"A B A",
			"A A A",
		}
		expected := []string{
			// lru: []
			"A A A",
			// lru: ["A A A"]
			"A A [*]",
			// lru: ["A A *"]
			"A B A",
			// lru: ["A B A", "A A *"]
			"A B [*]",
			// lru: ["A B *", "A A *"]
			"A C A",
			// lru: ["A C A", "A B *"]
			"A C [*]",
			// lru: ["A C *", "A B *"]
			"A B [*]",
			// Message "B A A" was normalized because the template "B A *" is
			// still present in the cache.
			// lru: ["B A *", "C A *"]
			"A A A",
			// Message "A A A" was not normalized because the template "C A A"
			// pushed out the template "A A *" from the cache.
			// lru: ["A A A", "C A *"]
		}
		drain := newDrain(withMaxClusters(2))
		for i, rawLog := range rawLogs {
			cluster, _ := drain.addLogMessage(rawLog)
			assert.Equal(t, expected[i], cluster.getTemplate())
		}
	})
}

func TestSeqDistance(t *testing.T) {
	testData := []struct {
		seq1          string
		seq2          string
		includeParams bool
		relVal        float32
		paramCount    int64
		err           error
	}{
		{
			seq1:          "abc 123 ooo",
			seq2:          "abc 456 ooo",
			includeParams: true,
			relVal:        0.6666667,
			paramCount:    0,
			err:           nil,
		},
		{
			seq1:          "456 123 ooo",
			seq2:          "abc 456 ooo",
			includeParams: true,
			relVal:        0.33333334,
			paramCount:    0,
			err:           nil,
		},
		{
			seq1:          "abc 123 [*]",
			seq2:          "abc 456 ooo",
			includeParams: true,
			relVal:        0.6666667,
			paramCount:    1,
			err:           nil,
		},
		{
			seq1:          "A A [*]",
			seq2:          "A B A",
			includeParams: false,
			relVal:        0.33333334,
			paramCount:    1,
			err:           nil,
		},
	}
	for _, data := range testData {
		t.Run(fmt.Sprintf("compare seq1: %s, seq2: %s, relVal: %v, paramCount: %v",
			data.seq1, data.seq2, data.relVal, data.paramCount), func(t *testing.T) {

			drain := newDrain()
			relVal, paramCount, _ := drain.getSeqDistance(
				utils.GetStringTokens(data.seq1), utils.GetStringTokens(data.seq2), data.includeParams)
			assert.Equal(t, data.relVal, relVal, "relVal not equals")
			assert.Equal(t, data.paramCount, paramCount, "paramCount not equals")
		})
	}
}

func TestFastMatch(t *testing.T) {
	testData := []struct {
		clusters    []*LogCluster
		tokens      []string
		bestMatchId int64
	}{
		{
			clusters: []*LogCluster{
				{
					id: 1,
					logTemplateTokens: []string{
						"abc", "123", "345",
					},
				},
				{
					id: 2,
					logTemplateTokens: []string{
						"abc", "bcd", "345",
					},
				},
			},
			tokens:      []string{"abc", "bcd", "123"},
			bestMatchId: 2,
		},
	}
	drain := newDrain()

	for _, data := range testData {
		for _, cluster := range data.clusters {
			drain.idToCluster.Add(cluster.id, nil)
		}
	}
	for i, data := range testData {
		t.Run(fmt.Sprintf("test %v group", i), func(t *testing.T) {
			bestMatchCluster := drain.fastMatch(data.clusters, data.tokens, drain.sim, true)
			assert.NotNil(t, bestMatchCluster)
			assert.Equal(t, data.bestMatchId, bestMatchCluster.id)
		})
	}
}

func TestAddSeqToPrefixTree(t *testing.T) {
	t.Run("add log template to empty tree", func(t *testing.T) {
		drain := newDrain(withDepth(5))
		rootNode := newRootTreeNode()
		logTemplate := []string{"abc", "aaa", "bcd", "def"}
		logCluster := newLogCluster(1, logTemplate)
		drain.addSeqToPrefixTree(rootNode, logCluster)
		assert.Equal(t, 1, len(rootNode.lengthChildren))
		assert.Equal(t, rootNode.lengthChildren[len(logTemplate)].length, len(logTemplate))
		lengthNode := rootNode.lengthChildren[len(logTemplate)]
		assert.NotNil(t, lengthNode.internalChildren[logTemplate[0]])
		firstInternalLayerNode := lengthNode.internalChildren[logTemplate[0]]
		assert.NotNil(t, firstInternalLayerNode.internalChildren[logTemplate[1]])
		secondInternalLayerNode := firstInternalLayerNode.internalChildren[logTemplate[1]]
		assert.Equal(t, 1, len(secondInternalLayerNode.clusters))
		assert.Equal(t, logCluster, secondInternalLayerNode.clusters[0])
	})
}

func TestMatch(t *testing.T) {
	t.Run("test match", func(t *testing.T) {
		drain := newDrain()
		drain.addLogMessage("aa aa aa")
		drain.addLogMessage("aa aa bb")
		drain.addLogMessage("aa aa cc")
		drain.addLogMessage("xx yy zz")
		c := drain.match("aa aa tt", SEARCH_STRATEGY_NEVER)
		assert.Equal(t, int64(1), c.id)
		c = drain.match("xx yy zz", SEARCH_STRATEGY_NEVER)
		assert.Equal(t, int64(2), c.id)
		c = drain.match("xx yy rr", SEARCH_STRATEGY_NEVER)
		assert.Nil(t, c)
		c = drain.match("nothing", SEARCH_STRATEGY_NEVER)
		assert.Nil(t, c)
	})
}

func TestDiff(t *testing.T) {
	t.Run("test diff", func(t *testing.T) {
		drain1 := newDrain()
		drain1.addLogMessage("aa bb cc")
		drain1.addLogMessage("aa bb dd")
		drain1.addLogMessage("aa bb ee")
		drain1.addLogMessage("bb cc dd ff")
		drain1.addLogMessage("bb cc dd gg")
		drain1.addLogMessage("bb cc dd ee kk")

		drain2 := newDrain()
		drain2.addLogMessage("aa bb cc")
		drain2.addLogMessage("aa bb dd")
		drain2.addLogMessage("bb cc dd ff")
		drain2.addLogMessage("bb cc dd ee")
		drain2.addLogMessage("bb cc dd gg")
		drain2.addLogMessage("bb cc dd ee kk")
		drain2.addLogMessage("bb cc dd ee kk ee")

		diff := drain2.diff(drain1)
		assert.Equal(t, 4, len(diff))

		assert.Equal(t, -1, diff[1].DiffNum)
		assert.Equal(t, float32(-0.33333334), diff[1].DiffRate)
		assert.Equal(t, LOG_CLUSTER_DIFF_TYPE_DECREASE, diff[1].DiffType)

		assert.Equal(t, 1, diff[2].DiffNum)
		assert.Equal(t, float32(0.5), diff[2].DiffRate)
		assert.Equal(t, LOG_CLUSTER_DIFF_TYPE_INCREASE, diff[2].DiffType)

		assert.Equal(t, 0, diff[3].DiffNum)
		assert.Equal(t, float32(0), diff[3].DiffRate)
		assert.Equal(t, LOG_CLUSTER_DIFF_TYPE_EQUAL, diff[3].DiffType)

		assert.Equal(t, 1, diff[4].DiffNum)
		assert.Equal(t, float32(0), diff[4].DiffRate)
		assert.Equal(t, LOG_CLUSTER_DIFF_TYPE_NEW, diff[4].DiffType)
	})
}
