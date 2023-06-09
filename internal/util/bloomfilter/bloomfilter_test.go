package bloomfilter

import (
	"math"
	"testing"
)

func Test(t *testing.T) {
	config := Config{MiscalRate: 0.00001, AddSize: 100, StorePath: "../../../data/bloomfilter"}

	nub := math.Log(2) * math.Log(2)
	t.Log(nub)
	bf := NewBloomFilter(config)
	var docid uint64
	docid = 123
	bf.AddNub(docid)
	ret := bf.CheckNub(docid)
	t.Log(ret)
	docid = 120
	bf.AddNub(docid)
	ret = bf.CheckNub(docid)
	t.Log(ret)
	docid = 124
	ret = bf.CheckNub(docid)
	t.Log(ret)
	bf.Save2File()

	bf = NewBloomFilter(config)
	mbSize := bf.Size() / 8 / 1000 / 1000
	t.Log("MbSize:", mbSize)
	ret = bf.CheckNub(123)
	t.Log(ret)
	ret = bf.CheckNub(120)
	t.Log(ret)
	ret = bf.CheckNub(124)
	t.Log(ret)
	DeleteBloomFile(config.StorePath)

	/*	bf = NewBloomFilter(0.00001, 100000000, "../../../data/bloomfilter")
		mbSize = bf.Size() / 8 / 1000 / 1000
		t.Log("MbSize:", mbSize)
		docid = 125
		bf.AddNub(docid)
		ret = bf.CheckNub(docid)
		t.Log(ret)*/
}
