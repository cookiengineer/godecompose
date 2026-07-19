arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Set { name: "reflect.Value.Set"; library: "go-stdlib"; instr match { CALL reflect___Value__Set } gen { $v.Set($x); } bind { v as "v"; x as "x"; } }
