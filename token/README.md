An example implementation of using a continuation token for API paging

Inspired by:
  * https://blog.philipphauer.de/web-api-pagination-timestamp-id-continuation-token
  * http://blog.vermorel.com/journal/2015/5/8/nearly-all-web-apis-get-paging-wrong.html

To run...

```
# one session
docker-compose up

# another
go test

```

Possible improvements:
  * create a client/server with headers to indicate continuation token
  * play around with isolation levels in postgres transactions to expose possible phantom or nonrepeatable read
