# TODO

## Pending Tasks

* wrap the POC in the cobra framework
* make 'datatype' a struct
* modify imagePull() to download files to a specified dir
* research why a request with no user agent gives a 403 forbidden
* research how to properly use 'goquerySelector' for more advanced HTMLElement queries https://github.com/PuerkitoBio/goquery
* there is a looping problem when downloading the images... figure out how to fix that. probably deals with the contexts or async
    * the collector should spawn a new goroutine for each image it's about to download, instead of doing it sequentially
    * the program works asynchronously with the c.Wait() calls, but the output is wonky... figure that out too
* Update the LICENSE information to make it as restrictive as possible


## Completed Tasks

* implement POC using a simple main func
* use the gocolly/colly library for scraping functionality - http://go-colly.org/ 