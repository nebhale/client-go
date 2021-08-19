# client-go <!-- omit in toc -->

[![Tests](https://github.com/nebhale/client-go/workflows/Tests/badge.svg?branch=main)](https://github.com/nebhale/client-go/actions/workflows/tests.yaml)
[![GoDoc](https://godoc.org/github.com/nebhale/client-go?status.svg)](https://godoc.org/github.com/nebhale/client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/nebhale/client-go)](https://goreportcard.com/report/github.com/nebhale/client-go)
[![codecov](https://codecov.io/gh/nebhale/client-go/branch/main/graph/badge.svg)](https://codecov.io/gh/nebhale/client-go)

`client-go` is a library to access [Service Binding Specification for Kubernetes](https://k8s-service-bindings.github.io/spec/) conformant Service Binding [Workload Projections](https://k8s-service-bindings.github.io/spec/#workload-projection).

## Example

```golang
import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/nebhale/client-go/bindings"
	"os"
)

func main() {
	b := bindings.FromServiceBindingRoot()
	b = bindings.Filter(b, "postgresql")
	if len(b) != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Incorrect number of PostgreSQL drivers: %d\n", len(b))
		os.Exit(1)
	}

	u, ok := bindings.Get(b[0], "url")
	if !ok {
		_, _ = fmt.Fprintln(os.Stderr, "No URL in binding")
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), u)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// ...
}
```

## License

Apache License v2.0: see [LICENSE](./LICENSE) for details.
