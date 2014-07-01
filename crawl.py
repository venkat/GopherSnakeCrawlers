#!/usr/bin/env python
# monkey-patch
import gevent.monkey
gevent.monkey.patch_all()

import sys
import hashlib
import re

import requests
from requests.exceptions import ConnectionError, MissingSchema
import gevent.pool
from gevent.queue import JoinableQueue

source = sys.argv[1] #source link
num_worker_threads = int(sys.argv[2]) #specifying how many workers fetch concurrently
num_to_crawl = int(sys.argv[3]) #maximum no. of pages to fetch

crawled = 0
links_added = 0
q = JoinableQueue() #JoinableQueue lets us wait till all the tasks in the queue are marked as done.

#This function does the actual work of fetching the link and 
#adding the extracted links from the page content into the queue
def do_work(link, crawler_id):
    global crawled, links_added

    #NOTE: uncomment this line to get extra details on what's happening
    #print 'crawling', crawled, crawler_id, link

    #Fetch the link
    try:
        response = requests.get(link) 
        response_content = response.content
    except (ConnectionError, MissingSchema):
        return
    crawled += 1

    #Some non-IO bound work on the content In the real world, there 
    #would be some heavy-duty parsing, DOM traversal here.
    
    m = hashlib.md5()
    m.update(response_content)
    m.digest()

    #Extract the links and add them to the queue. Using links_added
    #counter to keep track of links fetched. Possible race condition.
    for link in re.findall('<a href="(http.*?)"', response_content):
        if links_added < num_to_crawl:
            links_added += 1
            q.put(link) 

#Worker spawned by gevent. Continously gets links, works on them and marks
#them as done.
def worker(crawler_id):
    while True:
        item = q.get()
        try:
            do_work(item, crawler_id)
        finally:
            q.task_done()

#Spawning worker threads.
crawler_id = 0
for i in range(num_worker_threads):
    gevent.spawn(worker, crawler_id)
    crawler_id += 1 

q.put(source)
links_added += 1

q.join()  # block until all tasks are done
