package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed plugin.wasm
var wasmFile []byte

func main() {
	ctx := context.Background()
	cfg := wazero.NewRuntimeConfig()
	r := wazero.NewRuntimeWithConfig(ctx, cfg)

	builder := r.NewHostModuleBuilder("env")
	wasi_snapshot_preview1.NewFunctionExporter().ExportFunctions(builder)
	if _, err := builder.Instantiate(ctx); err != nil {
		panic(err)
	}

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.CompileModule(ctx, wasmFile)
	if err != nil {
		panic(err)
	}
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()
	modCfg := wazero.NewModuleConfig().
		WithStdin(stdinR).
		WithStdout(stdoutW).
		WithStderr(os.Stderr)

	go func() {
		if _, err := r.InstantiateModule(ctx, mod, modCfg); err != nil {
			panic(err)
		}
	}()

	fmt.Fprintf(os.Stderr, "write buffer\n")
	if _, err := stdinW.Write([]byte("hello world\n")); err != nil {
		panic(err)
	}
	stdinW.Close()
	reader := bufio.NewReader(stdoutR)
	for {
		content, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		fmt.Fprintf(os.Stderr, "content = %s. err = %v\n", content, err)
		break
	}
}
