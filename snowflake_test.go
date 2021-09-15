package snowflake

import (
	"sync"
	"testing"
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
//single instance multi goroutine
func TestSnowFlake_simg_GetId(t *testing.T) {
	sf1, err := New(2, 2, 2, 2, T38)
	if err != nil {
		t.Error(err)
		return
	}
	var wg sync.WaitGroup
	num := 2400000
	ids := make([]int64, num)
	for i := 0; i < num; i++ {
		wg.Add(1)
		ii := i
		go func() {
			id1 := sf1.GetId()
			ids[ii] = id1
			wg.Done()
		}()
	}
	wg.Wait()
	m:=make(map[int64]interface{})
	for _,v:=range ids{
		m[v]=nil
	}
	if len(m)!=len(ids){
		t.Errorf("TestSnowFlake_GetId() test error,want:%d,got:%d",len(ids),len(m))
	}

}

//multi instance multi goroutine
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
	m:=make(map[int64]interface{})
	for _,v:=range ids{
		m[v]=nil
	}
	if len(m)!=len(ids){
		t.Errorf("TestSnowFlake_GetId() test error,want:%d,got:%d",len(ids),len(m))
	}

}
