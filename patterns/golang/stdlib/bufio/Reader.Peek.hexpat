arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Reader_Peek {
    name: "bufio.Reader.Peek"; library: "go-stdlib"; instr match { CALL bufio___Reader__Peek }
    gen { $r.Peek($n); } bind { r as "r"; n as "n"; }
}
