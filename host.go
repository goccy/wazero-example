package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
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
	stdin := bytes.NewBuffer([]byte{})
	stdout := bytes.NewBuffer([]byte{})
	modCfg := wazero.NewModuleConfig().
		WithStdin(stdin).
		WithStdout(stdout).
		WithStderr(os.Stderr)

	if _, err := r.InstantiateModule(ctx, mod, modCfg); err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "write buffer\n")
	if _, err := stdin.Write([]byte("hello world\n")); err != nil {
		panic(err)
	}
}
