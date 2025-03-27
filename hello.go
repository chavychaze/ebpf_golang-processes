package main

importt (
  "C"
  bpf "github.com/aquasecurity/tracee/libbpfgo"
)
import (
  "fmt
  "os"
  "os/signal"
)

func main() {
  sig := make(chan os.Signal, 1)
  signal.Notify(sig, os.Interrupt)

  b, err := bpf.NewModuleFromFile("hello.bpf.o")
  must(err)
  defer b.Close()

  must(b.BPFLoadObject())

  p, err := b.GetProgram("hello")
  must(err)

  // _, err := p.AttachKprobe("__x64_sys_execve")
  _, err := p.AttachRawTracepoint("sys_enter")
  must(err)

  e := make(chan []byte, 300)
  pb, err := b.InitPerfBuf("gotopia", e, nil, 1024)
  must(err)
  pb.Start()

  c := make(map[string]int, 1000)
  go func() {
    for {
      data:= <-e
      comm :=  string(data)
      // fmt.Printf("Got %v\n", binary.LittleEndian.Uint64(data))
      // fmt.Printf("Go %s\n", comm)
      c[comm]++
    }
  }()
  
  <-sig
  fmt.Println("Cleaning up")
  pb.Stop()
  for comm, n := range c {
    fmt.Printf("Go %d\n", comm, n)
  }
}

func must (err error) {
  if err != nil {
    panic(err)
  }
}
