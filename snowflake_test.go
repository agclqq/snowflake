package snowflake

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func BenchmarkSnowFlake_GetId(b *testing.B) {
	sf, err := New(2, 2, 2, 2, T38)
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		sf.GetId()
	}

}

// single instance multi goroutine
func TestSnowFlake_simg_GetId(t *testing.T) {
	sf1, err := New(2, 2, 2, 2, T38)
	if err != nil {
		t.Error(err)
		return
	}
	var wg sync.WaitGroup
	num := 2400000
	ids := make([]int64, num)
	start := time.Now()
	for i := 0; i < num; i++ {
		wg.Add(1)
		ii := i
		go func() {
			ids[ii] = sf1.GetId()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("time:", time.Since(start))
	m := make(map[int64]interface{})
	for _, v := range ids {
		m[v] = nil
	}
	if len(m) != len(ids) {
		t.Errorf("TestSnowFlake_GetId() test error,want:%d,got:%d", len(ids), len(m))
	}

}

// multi instance multi goroutine
func TestSnowFlake_mimg_GetId(t *testing.T) {
	sf1, err := New(2, 2, 2, 2, T38)
	if err != nil {
		t.Error(err)
		return
	}
	sf2, err := New(1, 2, 2, 2, T38)
	if err != nil {
		t.Error(err)
		return
	}
	var wg sync.WaitGroup
	num := 1200000
	ids1 := make([]int64, num)
	ids2 := make([]int64, num)
	for i := 0; i < num; i++ {
		wg.Add(2)
		ii := i
		go func() {
			id1 := sf1.GetId()
			ids1[ii] = id1
			wg.Done()
		}()
		go func() {
			id2 := sf2.GetId()
			ids2[ii] = id2
			wg.Done()
		}()
	}
	wg.Wait()

	ids := append(ids1, ids2...)
	m := make(map[int64]interface{})
	for _, v := range ids {
		m[v] = nil
	}
	if len(m) != len(ids) {
		t.Errorf("TestSnowFlake_GetId() test error,want:%d,got:%d", len(ids), len(m))
	}
}

func randMm(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

var tryNum = 1000000

// Multi-terminal technique
func TestSnowFlake_GetRandomId(t *testing.T) {
	var wg sync.WaitGroup
	_dataCenterBit := 8
	_machineBit := 8
	_timeBit := int(T38)
	num := -1 ^ (-1 << (63 - _dataCenterBit - _machineBit - _timeBit))

	objNum := 500
	obj := make([]*SnowFlake, objNum)
	rs := make([][]int64, objNum)

	for i := 0; i < objNum; i++ {
		wg.Add(1)
		ii := i
		go func() {
			defer wg.Done()
			sf, err := New(randMm(0, -1^(-1<<8)), 8, randMm(0, -1^(-1<<8)), 8, T38)
			if err != nil {
				t.Error(err)
				return
			}
			obj[ii] = sf
		}()
	}
	wg.Wait()
	for i := 0; i < objNum; i++ {
		wg.Add(1)
		ii := i
		go func() {
			wg.Done()
			for j := 0; j < num; j++ {
				rs[ii] = append(rs[ii], obj[ii].GetId())
			}
		}()
	}
	wg.Wait()

	k := 0
	m := make(map[int64]interface{})
	for _, v := range rs {
		for _, vv := range v {
			k++
			m[vv] = nil
		}
	}

	if len(m) != k {
		fmt.Println(tryNum)
		t.Errorf("TestSnowFlake_GetRandomId() test error,want:%d,got:%d", k, len(m))
		return
	}

	tryNum--
	if tryNum > 0 {
		TestSnowFlake_GetRandomId(t)
	}
}
