// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

/*
 * Helpers for building cmd/go and cmd/cgo.
 */

// mkzdefaultcc writes zdefaultcc.go:
//
//	package main
//	const defaultCC = <defaultcc>
//	const defaultCXX = <defaultcxx>
//	const defaultPkgConfig = <defaultpkgconfig>
//
// It is invoked to write cmd/go/internal/cfg/zdefaultcc.go
// but we also write cmd/cgo/zdefaultcc.go
func mkzdefaultcc(dir, file string) {
	if strings.Contains(file, filepath.FromSlash("go/internal/cfg")) {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "// Code generated by go tool dist; DO NOT EDIT.\n")
		fmt.Fprintln(&buf)
		fmt.Fprintf(&buf, "package cfg\n")
		fmt.Fprintln(&buf)
		fmt.Fprintf(&buf, "const DefaultPkgConfig = `%s`\n", defaultpkgconfig)
		buf.WriteString(defaultCCFunc("DefaultCC", defaultcc))
		buf.WriteString(defaultCCFunc("DefaultLn", defaultld))
		buf.WriteString(defaultCCFunc("DefaultCXX", defaultcxx))
		buf.WriteString(defaultToolchainCCFunc("DefaultToolchainLd", defaultld))
		buf.WriteString(defaultToolchainCCFunc("DefaultToolchainAsm", defaultasm))
		buf.WriteString(defaultToolchainCCFunc("DefaultToolchainCC", defaultcc))
		buf.WriteString(defaultToolchainCCFunc("DefaultToolchainCXX", defaultcxx))
		writefile(buf.String(), file, writeSkipSame)
		return
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// Code generated by go tool dist; DO NOT EDIT.\n")
	fmt.Fprintln(&buf)
	fmt.Fprintf(&buf, "package main\n")
	fmt.Fprintln(&buf)
	fmt.Fprintf(&buf, "const defaultPkgConfig = `%s`\n", defaultpkgconfig)
	buf.WriteString(defaultCCFunc("defaultCC", defaultcc))
	buf.WriteString(defaultCCFunc("defaultLd", defaultcc))
	buf.WriteString(defaultCCFunc("defaultCXX", defaultcxx))
	buf.WriteString(defaultToolchainCCFunc("defaultToolchainLd", defaultld))
	buf.WriteString(defaultToolchainCCFunc("defaultToolchainAsm", defaultasm))
	buf.WriteString(defaultToolchainCCFunc("defaultToolchainCC", defaultcc))
	buf.WriteString(defaultToolchainCCFunc("defaultToolchainCXX", defaultcxx))
	writefile(buf.String(), file, writeSkipSame)
}

func defaultCCFunc(name string, defaultcc map[string]string) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "func %s(goos, goarch string) string {\n", name)
	fmt.Fprintf(&buf, "\tswitch goos+`/`+goarch {\n")
	var keys []string
	for k := range defaultcc {
		if k != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&buf, "\tcase %q:\n\t\treturn %q\n", k, defaultcc[k])
	}
	fmt.Fprintf(&buf, "\t}\n")
	fmt.Fprintf(&buf, "\treturn %q\n", defaultcc[""])
	fmt.Fprintf(&buf, "}\n")

	return buf.String()
}

func defaultToolchainCCFunc(name string, defaultcc map[string]string) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "func %s(goos, goarch, toolchain string) string {\n", name)
	fmt.Fprintf(&buf, "\tswitch goos+`/`+goarch+`/`+toolchain {\n")
	var keys []string
	for k := range defaultcc {
		if k != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&buf, "\tcase %q:\n\t\treturn %q\n", k, defaultcc[k])
	}
	fmt.Fprintf(&buf, "\t}\n")
	fmt.Fprintf(&buf, "\t\treturn %q\n", defaultcc[""])
	fmt.Fprintf(&buf, "}\n")

	return buf.String()
}

// mkzcgo writes zosarch.go for cmd/go.
func mkzosarch(dir, file string) {
	// sort for deterministic zosarch.go file
	var list []string
	for plat := range cgoEnabled {
		list = append(list, plat)
	}
	sort.Strings(list)

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// Code generated by go tool dist; DO NOT EDIT.\n\n")
	fmt.Fprintf(&buf, "package cfg\n\n")
	fmt.Fprintf(&buf, "var OSArchSupportsCgo = map[string]bool{\n")
	for _, plat := range list {
		fmt.Fprintf(&buf, "\t%q: %v,\n", plat, cgoEnabled[plat])
	}
	fmt.Fprintf(&buf, "}\n")

	writefile(buf.String(), file, writeSkipSame)
}

// mkzcgo writes zcgo.go for the go/build package:
//
//	package build
//  var cgoEnabled = map[string]bool{}
//
// It is invoked to write go/build/zcgo.go.
func mkzcgo(dir, file string) {
	// sort for deterministic zcgo.go file
	var list []string
	for plat, hasCgo := range cgoEnabled {
		if hasCgo {
			list = append(list, plat)
		}
	}
	sort.Strings(list)

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// Code generated by go tool dist; DO NOT EDIT.\n")
	fmt.Fprintln(&buf)
	fmt.Fprintf(&buf, "package build\n")
	fmt.Fprintln(&buf)
	fmt.Fprintf(&buf, "const defaultCGO_ENABLED = %q\n", os.Getenv("CGO_ENABLED"))
	fmt.Fprintln(&buf)
	fmt.Fprintf(&buf, "var cgoEnabled = map[string]bool{\n")
	for _, plat := range list {
		fmt.Fprintf(&buf, "\t%q: true,\n", plat)
	}
	fmt.Fprintf(&buf, "}\n")

	writefile(buf.String(), file, writeSkipSame)
}
