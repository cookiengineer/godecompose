arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_MapIndex { name: "reflect.Value.MapIndex"; library: "go-stdlib"; instr match { CALL reflect___Value__MapIndex } gen { $v.MapIndex($key); } bind { v as "v"; key as "key"; } }
