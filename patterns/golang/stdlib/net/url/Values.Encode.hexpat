arch x86_64;
platform linux, darwin, windows, freebsd;

pattern url_Values_Encode {
    name: "url.Values.Encode"; library: "go-stdlib";
    instr match { CALL net_url___Values__Encode }
    gen { $v.Encode(); } bind { v as "v"; }
}
