arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_NewScanner {
    name: "bufio.NewScanner"; library: "go-stdlib"; instr match { CALL bufio_NewScanner }
    gen { bufio.NewScanner($r); } bind { r as "r"; }
}
