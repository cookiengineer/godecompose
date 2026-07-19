arch x86_64;
platform linux, darwin, windows, freebsd;

pattern url_Values_Set {
    name: "url.Values.Set"; library: "go-stdlib";
    instr match { CALL net_url___Values__Set }
    gen { $v.Set($key, $value); } bind { v as "v"; key as "key"; value as "value"; }
}
