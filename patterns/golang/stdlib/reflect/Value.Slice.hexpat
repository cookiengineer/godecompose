arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Slice { name: "reflect.Value.Slice"; library: "go-stdlib"; instr match { CALL reflect___Value__Slice } gen { $v.Slice($i, $j); } bind { v as "v"; i as "i"; j as "j"; } }
