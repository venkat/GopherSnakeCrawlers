GopherSnakeCrawlers
===================

Comparing simple Go and Python web crawler

For each crawler, the first argument is the source link, the second is the number of workers and the third is the number of pages to fetch.

For example: `python crawl.py http://google.com 1 10`

`run.sh` runs the crawler multiple times with different worker, page fetch parameters. The first argument is to specify go or python and the second argument is the source link. For go, it expects an executable with the name `crawl` be present in the current directory.

For example: `./run.sh python http://google.com`

More details at: http://venkat.io/posts/concurrent-crawling/
