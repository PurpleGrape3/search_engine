package rocksdbwrapper

/*
import (
	"errors"
	"search_engine/internal/util/log"
	"search_engine/internal/util/tools"
	"time"

	"github.com/tecbot/gorocksdb"
)

type rocksdbWrapper struct {
	DB *gorocksdb.DB
	RO *gorocksdb.ReadOptions
	WO *gorocksdb.WriteOptions
}

func NewRocksdb(path string) (*rocksdbWrapper, error) {
	rs := new(rocksdbWrapper)
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	var err error
	rs.DB, err = gorocksdb.OpenDb(opts, path)
	rs.RO = gorocksdb.NewDefaultReadOptions()
	rs.WO = gorocksdb.NewDefaultWriteOptions()
	return rs, err
}

func (rs *rocksdbWrapper) Set(k []byte, v []byte) error {
	defer func(cost func() time.Duration) {
		log.Debugf("trackid:%v, cost: %.3f ms", 0, float64(cost().Microseconds())/1000.0)
	}(tools.TimeCost())

	log.Debugf("key:%s", tools.Bytes2Str(k))
	return rs.DB.Put(rs.WO, k, v)
}

func (rs *rocksdbWrapper) Get(k []byte) ([]byte, error) {
	defer func(cost func() time.Duration) {
		log.Debugf("trackid:%v, cost: %.3f ms", 0, float64(cost().Microseconds())/1000.0)
	}(tools.TimeCost())

	s, err := rs.DB.Get(rs.RO, k)
	v := s.Data()
	log.Debugf("key:%s", tools.Bytes2Str(k))
	if len(v) == 0 {
		return nil, errors.New("data empty")
	}
	return v, err
}

func (rs *rocksdbWrapper) Delete(k []byte) error {
	return rs.DB.Delete(rs.WO, k)
}

func (rs *rocksdbWrapper) Close() error {
	rs.DB.Close()
	return nil
}*/
