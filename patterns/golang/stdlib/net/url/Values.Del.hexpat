arch x86_64;
platform linux, darwin, windows, freebsd;

pattern url_Values_Del {
    name: "url.Values.Del"; library: "go-stdlib";
    instr match { CALL net_url___Values__Del }
    gen { $v.Del($key); } bind { v as "v"; key as "key"; }
}
