arch x86_64;
platform linux, darwin, windows, freebsd;

pattern unsafe_Sizeof { name: "unsafe.Sizeof"; library: "go-stdlib"; instr match { CALL unsafe_Sizeof } gen { unsafe.Sizeof($x); } bind { x as "x"; } }
pattern unsafe_Offsetof { name: "unsafe.Offsetof"; library: "go-stdlib"; instr match { CALL unsafe_Offsetof } gen { unsafe.Offsetof($x); } bind { x as "x"; } }
pattern unsafe_Alignof { name: "unsafe.Alignof"; library: "go-stdlib"; instr match { CALL unsafe_Alignof } gen { unsafe.Alignof($x); } bind { x as "x"; } }
pattern unsafe_Pointer { name: "unsafe.Pointer"; library: "go-stdlib"; instr match { CALL unsafe_Pointer } gen { unsafe.Pointer($p); } bind { p as "p"; } }
pattern unsafe_Add { name: "unsafe.Add"; library: "go-stdlib"; instr match { CALL unsafe_Add } gen { unsafe.Add($p, $len); } bind { p as "p"; len as "len"; } }
pattern unsafe_Slice { name: "unsafe.Slice"; library: "go-stdlib"; instr match { CALL unsafe_Slice } gen { unsafe.Slice($p, $len); } bind { p as "p"; len as "len"; } }
pattern unsafe_SliceData { name: "unsafe.SliceData"; library: "go-stdlib"; instr match { CALL unsafe_SliceData } gen { unsafe.SliceData($s); } bind { s as "s"; } }
pattern unsafe_String { name: "unsafe.String"; library: "go-stdlib"; instr match { CALL unsafe_String } gen { unsafe.String($p, $len); } bind { p as "p"; len as "len"; } }
pattern unsafe_StringData { name: "unsafe.StringData"; library: "go-stdlib"; instr match { CALL unsafe_StringData } gen { unsafe.StringData($s); } bind { s as "s"; } }
