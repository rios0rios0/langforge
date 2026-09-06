// Package toolchain holds the build validator and runtime manager that the
// language providers compose. Validating a repository means running the
// ecosystem's lint commands and then its build commands, and managing its
// runtime means describing an SDK and asking its binary which version is
// installed. Neither differs between ecosystems except in the commands and
// names involved, so each language package supplies only those and shares the
// implementation kept here, the way javagradle and javamaven share the java
// package.
package toolchain
