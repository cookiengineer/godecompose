arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_NewWriter {
    name: "bufio.NewWriter"; library: "go-stdlib"; instr match { CALL bufio_NewWriter }
    gen { bufio.NewWriter($wr); } bind { wr as "wr"; }
}
