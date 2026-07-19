arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Writer_Reset {
    name: "bufio.Writer.Reset"; library: "go-stdlib"; instr match { CALL bufio___Writer__Reset }
    gen { $w.Reset($wr); } bind { w as "w"; wr as "wr"; }
}
