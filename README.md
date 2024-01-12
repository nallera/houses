# Houses image downloader

## Solution notes

I designed the solution based on a hexagonal architecture, but without the server and endpoints since it is a script.
The idea is to run the tasks that depend on external services concurrently, to optimize the execution time. A retry strategy was implemented as well, to take into account the 'flakiness' of the web sources.

## Run parameters

- _houses_: is the total number of houses (and images) to retrieve
- _pages_: the number of pages in which to distribute the number of houses

Note: the number of houses must be greater than or equal to the number of pages

### Run example

`go run main.go -houses 10 -pages 4`

This will retrieve 12 houses, in 4 pages (3 houses per page), and will take into account only the first 10.

## Tests 

Some unit testing was included, some _TODO_ sections as well. 

### Run example

`go test -v ./...`