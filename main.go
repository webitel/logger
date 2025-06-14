package main

import (
	"fmt"
	"github.com/webitel/logger/cmd/logger"
	// -------------------- plugin(s) -------------------- //
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/log/otlp"
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/log/stdout"
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/metric/otlp"
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/metric/stdout"
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/trace/otlp"
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/trace/stdout"
)

//go:generate go run github.com/bufbuild/buf/cmd/buf@latest generate --template buf/buf.gen.yaml
//go:generate go run github.com/bufbuild/buf/cmd/buf@latest generate --template buf/buf.gen.storage.yaml
//go:generate go run github.com/bufbuild/buf/cmd/buf@latest generate --template buf/buf.gen.webitel-go.yaml
//go:generate go run github.com/bufbuild/buf/cmd/buf@latest generate --template buf/buf.gen.engine.yaml

func main() {

	err := logger.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
