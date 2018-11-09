## Feedback Agent

The loadbalancer.org feedback agent is used to dynamically set real service weight depending on it's available system resources in <a href="http://www.haproxy.org/">HAProxy<a/>.
  
For an in-depth look into how our feedback agent works please refer to this blog post: <a href="http://www.loadbalancer.org/blog/open-source-windows-service-for-reporting-server-load-back-to-haproxy-load-balancer-feedback-agent/">Open Source Windows service for reporting server load back to HAProxy (load balancer feedback agent).</a>
  
## Features

- CPU metric
- RAM metric
- TCP connections metric
- Read/reload from config
- Halt/Down/Normal Status States

## Prerequisites

* Go v1.9 or later
* Windows

## Build
Please follow these instructions to build the feedback agent:

```
go get -d
```

```
go build
```
## Support
If you require assistence with our feedback agent please contact us at support@loadbalancer.org

## License
GNU General Public License, version 2

