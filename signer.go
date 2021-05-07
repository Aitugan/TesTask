// сюда писать код
package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for _, eachJob := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		go func(j job, in, out chan interface{}) {
			j(in, out)
			close(out)
			wg.Done()
		}(eachJob, in, out)
		in = out
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for i := range in {
		hash := strconv.Itoa(i.(int))
		md5 := DataSignerMd5(hash)
		wg.Add(1)
		go singleHashWorker(wg, hash, md5, out)
	}
	wg.Wait()
}

func singleHashWorker(wg *sync.WaitGroup, data string, md5 string, out chan interface{}) {
	crc32Ch := make(chan string)
	go func() {
		crc32Ch <- DataSignerCrc32(data)
	}()

	md5Ch := make(chan string)
	go func() {
		md5Ch <- DataSignerCrc32(md5)
	}()
	defer wg.Done()

	hash := <-crc32Ch + "~" + <-md5Ch
	out <- hash

}

// func MultiHash(in, out chan interface{}) {
// 	wg := &sync.WaitGroup{}
// 	for i := range in {
// 		wg.Add(1)
// 		go func(i, out chan interface{}) {
// 			defer wg.Done()
// 			hash := make([]string, 6)
// 			for th := 0; th < 6; th++ {
// 				data := fmt.Sprintf("%v%v", th, h)
// 				crc32hash := DataSignerCrc32(data)
// 				hash[th] = crc32hash

// 			}
// 		}(i, out)
// 	}

// 	wg.Wait()
// }

func CombineResults(in, out chan interface{}) {
	var results []string

	for i := range in {
		results = append(results, i.(string))
	}

	sort.Strings(results)
	out <- strings.Join(results, "_")
}
