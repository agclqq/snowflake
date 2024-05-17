package snowflake

import (
	"errors"
	"sync"
	"time"
)

const idBit = 63

type timeBit int64

const (
	T38 timeBit = iota + 38 //When the time bit is 38, the time range is 8.71 years, or about 3181.4 days
	T39                     //When the time bit is 39, the time range is 17.43 years, or about 6362.9 days
	T4O                     //When the time bit is 40, the time range is 34.86 years, or about 12,725.8 days
	T41                     //When the time bit is 41, the time range is 69.73 years, or about 25451.6 days
	T42                     //When the time bit is 42, the time range is 139.4 years, or about 50903.3 days
)

var lastTime int64
var sequence int64
var mx sync.Mutex               //To solve multi-instance concurrency
const startTime = 1630425600000 // 2021-09-01 00:00:00.000

type SnowFlake struct {
	dataCenter    int64
	dataCenterBit int64
	machine       int64
	machineBit    int64
	tb            timeBit
}

func New(dataCenter, dataCenterBit, machine, machineBit int64, tb timeBit) (*SnowFlake, error) {
	if dataCenterBit+machineBit+int64(tb) > idBit {
		return nil, errors.New("overlong digit,it's best to leave a few bits for the sequence")
	}
	if dataCenter > -1^(-1<<dataCenterBit) {
		return nil, errors.New("dataCenter overlong")
	}
	if machine > -1^(-1<<machineBit) {
		return nil, errors.New("machine overlong")
	}
	if getNowMil() < startTime {
		return nil, errors.New("the startTime must be less than now")
	}
	return &SnowFlake{
		dataCenter:    dataCenter,
		dataCenterBit: dataCenterBit,
		machine:       machine,
		machineBit:    machineBit,
		tb:            tb,
	}, nil
}

func (sf *SnowFlake) GetId() int64 {
	dcLeftMove := idBit - sf.dataCenterBit
	mcLeftMove := dcLeftMove - sf.machineBit
	tnLeftMove := mcLeftMove - int64(sf.tb)
	sequenceBit := tnLeftMove
	var now int64
	var id int64
	mx.Lock()
	for true {
		now = getNowMil()
		if lastTime == 0 || now == lastTime {
			if lastTime == 0 {
				lastTime = now
			}
			if sequence > -1^(-1<<sequenceBit) {
				continue // If the sequence value increases to the maximum value, wait for the next moment
			}
			id = sf.dataCenter<<dcLeftMove | sf.machine<<mcLeftMove | (now-startTime)<<tnLeftMove | sequence
			sequence++
			break
		} else {
			sequence = 0
			lastTime = now
		}
	}
	mx.Unlock()
	return id
}

func getNowMil() int64 {
	return time.Now().UnixNano() / 1e6
}
