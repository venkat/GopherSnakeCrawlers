package main

//The builtins are limited. Making a lot of imports necessary
import ("sync"
        "net/http"
        "regexp"
        "io/ioutil"
        "os"
        "bytes"
        "fmt"
        "strconv"
        "runtime"
        "crypto/md5"
        "io"
)

var source = os.Args[1] //source link
var num_worker_threads, _ = strconv.Atoi(os.Args[2]) //specifying how many workers fetch concurrently
var num_to_crawl, _ = strconv.Atoi(os.Args[3]) //maximum no. of pages to fetch

var crawled = make(chan int, num_to_crawl) //using a buffered channel to count page fetches
var links = make(chan string, num_to_crawl) //using a buffered channels as a queue of links

func do_work(link string, crawler_id int) {
    //fmt.Println("crawling", crawler_id, link)
    re := regexp.MustCompile(`<a href="(http.*?)"`)
    resp, err := http.Get(link)
    if err != nil {
        return
    }
    defer resp.Body.Close()
    content, _ := ioutil.ReadAll(resp.Body)
    contentString := bytes.NewBuffer(content).String()
    h := md5.New()
    io.WriteString(h, contentString)
    var _ = h.Sum(nil)

    //Try to add a link to the queue of links. If it is full, the default case
    //returns as there is no point in adding more links to the queue as our
    //maximum page fetches is limited anyways.
    for _, match := range re.FindAllStringSubmatch(contentString, -1) {
        select {
        case links <- match[1]:
        default:
            return
        }
    }
}

func worker(crawler_id int) {
    //If the crawled channel's buffer is full, no more pages to fetch
    //so no more work to do.
    for {
        select {
        case crawled <- 1:
            do_work(<-links, crawler_id)
        default:
            return
        }
    }
}

func main() {
    var _ = fmt.Println
    //Try to make the workers use all the logical CPUs in the machine.
    runtime.GOMAXPROCS(runtime.NumCPU())
    var wg sync.WaitGroup
    links <- source

    for i:=0; i < num_worker_threads; i++ {
        // Increment the WaitGroup counter.
        wg.Add(1)
        // Launch a goroutine worker.
        go func(crawler_id int) {
                // Decrement the counter when the goroutine completes.
                defer wg.Done()
                worker(crawler_id)
        }(i)
    }
    // Wait for all the workers to finish.
    wg.Wait()
    close(crawled)
    close(links)
}
