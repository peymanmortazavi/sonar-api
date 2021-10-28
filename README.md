# sonar-api
Sonar API Query

# How to Run

In order to run in docker simply run:

```bash
$ make run-image
```

This will start listening on port 6000, using the CSV file stored in `example/pings.csv`

You can use the following query parameters:

`after` and `before` that are UNIX time stamps

`depth` which is a number

`region` which is a commma separated list of numbers representing lon,lat of northwest and southwest

