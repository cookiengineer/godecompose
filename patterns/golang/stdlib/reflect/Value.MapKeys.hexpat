arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_MapKeys { name: "reflect.Value.MapKeys"; library: "go-stdlib"; instr match { CALL reflect___Value__MapKeys } gen { $v.MapKeys(); } bind { v as "v"; } }
