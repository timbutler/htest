# HTest
[HTest] - HTTP Testing Tool
-------

This is a very simple tool designed to be a HTTP endpoint for load balancer testing. By default, the htttrack listens on port 8000 and serves a very basic HTML page.

On an i5-6500, this is enough to serve over 130,000 requests per second so it should be sufficient for most large scale systems.

## Usage

You can get the full command line flags by running `htest --help`.

### Simple Instance
In it's most basic form, you can run a single instance:

    htest

This will run on the default port of 8000 on 127.0.0.1. If you need to run on a different port or IP, 
you can use:

    htest -port 8001 -ipaddr 192.168.10.50

This will then run on port `8001` and bind to the IP address `192.168.10.50`.     

### Multiple Instances
-----
For effectively testing load balancers, you will need to run multiple instances. Depending on your 
requirements, these could be on a single system or multiple systems. Here's a quick example of spinning up 5 instances:

    htest -port 8001 &
    htest -port 8002 &
    htest -port 8003 &
    htest -port 8004 &
    htest -port 8005 &

-----

### Varying Data

If you want to slow down some of the responses, you can do this with a simple http call via `curl`.

There are three different areas this can be changed:
 

| Field         | Value                                 | Example |
|---------------|---------------------------------------|---------|
| responsedelay | Delay Time (in ms)                    | 10      |
| failurerate   | Failure rate of requests (in percent) | 2       |
| jitter        | Variance of the results (in percent)  | 5       |
|

For example: 

    curl http://localhost:8000/svar/?responsedelay=10

This would set a 10 millisecond delay to the responses, so you can simulate the average response times from your server as well.
 To make it more realistic, you can add some `jitter` to the results so that they're not all exactly delayed by 10ms.

The `failurerate` then means that a certain percentage will return 503 errors instead of a 200 error. If you're wanting to test health monitoring of your system, 
this is an easy way to check for the server to be blacklisted by your loadbalancer. Again, the `jitter` figure comes into play here 
so that you have some variance in the results.

## Raw benchmark results

These are the basic results on a simple i5-6500 based workstation:

Run using Apache Benchmark: `ab -k -c 40 -n 500000 http://localhost:8000/`

**Results** 

    Requests per second:    139060.71 [#/sec] (mean)`
    Time per request:       0.288 [ms] (mean)
    Time per request:       0.007 [ms] (mean, across all concurrent requests)

The intent isn't to produce the world's fastest http server, but to simply ensure it's fast enough to run multiple instances behind a load balancer (or similar) and ensure it's not a limiting factor.

## TODO
 - Basic sanitisation of the variables
 - Add Binary releases + Travis CI
 - Create Docker instance

## License 

MIT License - https://en.wikipedia.org/wiki/MIT_License

## Contributing

Please feel free to report issues or submit pull requests for changes. 

The aim is to keep this program very simple and lightweight, so complex changes or 
large deviations in behaviour may not be accepted. As it's under the MIT license,
you can certainly fork and use as you see fit.

## Alternatives and other tools

- https://github.com/shopify/toxiproxy 