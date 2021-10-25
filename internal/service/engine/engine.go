package engine

import (
	"context"
	dm "search_engine/internal/model/datamanager"
	al "search_engine/internal/service/analyzer"
	"search_engine/internal/service/objs"
	rk "search_engine/internal/service/rank"
	bf "search_engine/internal/util/bloomfilter"
	"search_engine/internal/util/log"
	"sort"
	"strings"
)

type engine struct {
	docid       int64
	analyzer    al.Analyzer
	ranker      rk.Ranker
	bloomfilter *bf.BloomFilter
	datamanager *dm.Manager
}

func newEngine(analyzerStopWordPath string, dbPath string, dbHost string, dbPort string, dbPassword string, dbIndex int, dbTimeout int, bloomfilterMiscalRate float64, bloomfilterAddSize uint64) *engine {
	eg := new(engine)
	eg.docid = 0
	eg.analyzer = al.AnalyzerFactory(analyzerStopWordPath)
	eg.ranker = rk.RankerFactory()
	eg.bloomfilter = bf.NewBloomFilter(bloomfilterMiscalRate, bloomfilterAddSize)
	eg.datamanager = dm.NewManager(dbPath, dbHost, dbPort, dbPassword, dbIndex, dbTimeout)
	return eg
}

func (eg *engine) close() {
	eg.datamanager.Close()
}

func (eg *engine) retrieveDoc(ctx context.Context, retreiveTerms []objs.RetreiveTerm) objs.RecallPostingList {
	replUnion := make(objs.RecallPostingList, 0)
	replInters := make([]objs.RecallPostingList, 0)
	hasInter := false
	//TODO:开协程并发请求
	for _, terminfo := range retreiveTerms {
		if repl, err := eg.datamanager.Retrieve(ctx, terminfo.Field, terminfo.FieldData); err == nil {
			if terminfo.Operator == objs.Union {
				replUnion = append(replUnion, repl...)
			} else if terminfo.Operator == objs.Inter {
				replInters = append(replInters, repl)
				hasInter = true
			}
		}
	}

	//并集去重过滤
	sort.Sort(replUnion)
	replUniqUnion := make(objs.RecallPostingList, 0)
	docidSet := make(map[uint64]bool)
	for _, reposting := range replUnion {
		if !eg.filter(reposting, titleMust, priceStart, priceEnd) {
			if _, ok := docidSet[reposting.Docid]; !ok {
				docidSet[reposting.Docid] = true
				replUniqUnion = append(replUniqUnion, reposting)
			}
		}
	}

	if !hasInter {
		log.Debugf("trackid:%d, replUniqUnion:%v", ctx.Value("trackid").(uint64), replUniqUnion)
		return eg.ranker.Rank(replUniqUnion)
	}

	//交集去重过滤
	replUniqInters := make([]objs.RecallPostingList, 0)
	for _, repl := range replInters {
		sort.Sort(repl)
		plUniqInter := make(objs.RecallPostingList, 0)
		docidSet = make(map[uint64]bool)
		for _, reposting := range repl {
			if !eg.filter(reposting, titleMust, priceStart, priceEnd) {
				if _, ok := docidSet[reposting.Docid]; !ok {
					docidSet[reposting.Docid] = true
					plUniqInter = append(plUniqInter, reposting)
				}
			}
		}
		replUniqInters = append(replUniqInters, plUniqInter)
	}

	replCal := eg.calInter(replUniqUnion, replUniqInters)
	log.Debugf("trackid:%d, replUniqUnion:%v, replUniqInters:%v replCal:%v", ctx.Value("trackid").(uint64), replUniqUnion, replUniqInters, replCal)
	return eg.ranker.Rank(replCal)
}

//TODO：抽离成公共组件
//指针求交
func (eg *engine) calInter(replUniqUnion objs.RecallPostingList, replUniqInters []objs.RecallPostingList) objs.RecallPostingList {
	replUniqInters = append(replUniqInters, replUniqUnion)
	replsEnd := make([]int, len(replUniqInters))
	minEnd := len(replUniqInters[0])
	minIndex := 0
	for i, pl := range replUniqInters {
		replsEnd[i] = len(pl)
		if replsEnd[i] < minEnd {
			minEnd = replsEnd[i]
			minIndex = i
		}
	}
	repl := replUniqInters[minIndex]

	replCal := make(objs.RecallPostingList, 0)
	replUniqInters = append(replUniqInters[:minIndex], replUniqInters[minIndex+1:]...)
	replsEnd = append(replsEnd[:minIndex], replsEnd[minIndex+1:]...)
	replsStart := make([]int, len(replUniqInters))
	midBreak := false
	for _, reposting := range repl {
	reloop:
		for i := 0; i < len(replUniqInters); i++ {
			for {
				if replUniqInters[i][replsStart[i]].Docid < reposting.Docid {
					replsStart[i]++
					if replsStart[i] < replsEnd[i] {
						continue
					} else {
						goto finally
					}
				} else if replUniqInters[i][replsStart[i]].Docid == reposting.Docid {
					replsStart[i]++
					break
				} else {
					midBreak = true
					break reloop
				}
			}
		}
		if !midBreak {
			replCal = append(replCal, reposting)
		}
		midBreak = false
	}

finally:
	return replCal
}

func (eg *engine) filter(repo objs.RecallPosting, titleMust string, priceStart float64, priceEnd float64) bool {
	if !eg.docIsDel(repo.Docid) {
		if strings.Contains(repo.Title, titleMust) {
			if repo.Price >= priceStart && repo.Price <= priceEnd {
				return false
			}
		}
	}
	return true
}

func (eg *engine) addDoc(ctx context.Context, doc objs.Doc, docid uint64) {
	ps := eg.analyzer.Analysis(docid, doc)
	eg.datamanager.AddDoc(doc, docid, ps)
	log.Debugf("trackid:%d, docid:%d, ps:%v", ctx.Value("trackid").(uint64), docid, ps)
}

func (eg *engine) delDoc(docid uint64) {
	eg.bloomfilter.AddNub(docid)
}

func (eg *engine) docIsDel(docid uint64) bool {
	deleted := eg.bloomfilter.CheckNub(docid)
	return deleted
}
