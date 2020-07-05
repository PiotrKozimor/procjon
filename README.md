# procjon

[![Go Report Card](https://goreportcard.com/badge/github.com/PiotrKozimor/procjon?style=flat-square)](https://goreportcard.com/report/github.com/PiotrKozimor/procjon)
[![Coverage](https://codecov.io/gh/PiotrKozimor/procjon/branch/master/graph/badge.svg)](https://codecov.io/gh/PiotrKozimor/procjon)
[![Build Status Travis](https://img.shields.io/travis/PiotrKozimor/procjon.svg?style=flat-square&&branch=master)](https://travis-ci.org/github/PiotrKozimor/procjon)
[![LICENSE](https://img.shields.io/github/license/PiotrKozimor/procjon.svg?style=flat-square)](https://github.com/etcd-io/etcd/blob/master/LICENSE)
<!-- [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/etcd-io/etcd) -->
<!-- [![Releases](https://img.shields.io/github/release/PiotrKozimor/procjon.svg?style=flat-square)](https://github.com/etcd-io/etcd/releases) -->


**Note**: The `master` branch may be in an *unstable or even broken state* during development. Please use [releases][github-release] instead of the `master` branch in order to get stable binaries.

Procjon is simple and monitoring tool written in Go. Procjon is a server which collects statuses from procjonagents. When status sent by procjonagents changes, procjon sends **status** update to Slack. When **timeout** or error occurs, procjon sends **availability** update.

Procjonagents, e.g. procjonsystemd or procjonelastic communicates with procjon with gRPC interface to register **service** and update it's **status**:


![arch](doc/arch.svg)

```
service Procjon{
    rpc RegisterService(Service) returns (Empty) {}
    rpc SendServiceStatus(stream ServiceStatus) returns (Empty) {}
}
```

When registering service, human readable **statuses** are posted. **Timeout** is also configured:

```
message Service {
    string serviceIdentifier = 1;
    map<int32,string> statuses = 2;
    int32 timeout = 3;
}
```

Then, only **statusCode** update is sent to procjon:

```
message ServiceStatus {
    string serviceIdentifier = 1;
    int32 statusCode = 2;
}
```

Important operating principle of procjon is that reliability of monitoring server is higher that infrastructure to monitor. HA is not yet planned. Procjon was designed to deal with unreliable internal infrastructure and many processes which were set-up and then left forgotten (e.g. long term tests).

