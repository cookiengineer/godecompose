arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Abs { name: "math.Abs"; library: "go-stdlib"; instr match { CALL math_Abs } gen { math.Abs($x); } bind { x as "x"; } }
pattern math_Max { name: "math.Max"; library: "go-stdlib"; instr match { CALL math_Max } gen { math.Max($x, $y); } bind { x as "x"; y as "y"; } }
pattern math_Min { name: "math.Min"; library: "go-stdlib"; instr match { CALL math_Min } gen { math.Min($x, $y); } bind { x as "x"; y as "y"; } }
pattern math_Sqrt { name: "math.Sqrt"; library: "go-stdlib"; instr match { CALL math_Sqrt } gen { math.Sqrt($x); } bind { x as "x"; } }
pattern math_Pow { name: "math.Pow"; library: "go-stdlib"; instr match { CALL math_Pow } gen { math.Pow($x, $y); } bind { x as "x"; y as "y"; } }
pattern math_Sin { name: "math.Sin"; library: "go-stdlib"; instr match { CALL math_Sin } gen { math.Sin($x); } bind { x as "x"; } }
pattern math_Cos { name: "math.Cos"; library: "go-stdlib"; instr match { CALL math_Cos } gen { math.Cos($x); } bind { x as "x"; } }
pattern math_Tan { name: "math.Tan"; library: "go-stdlib"; instr match { CALL math_Tan } gen { math.Tan($x); } bind { x as "x"; } }
pattern math_Floor { name: "math.Floor"; library: "go-stdlib"; instr match { CALL math_Floor } gen { math.Floor($x); } bind { x as "x"; } }
pattern math_Ceil { name: "math.Ceil"; library: "go-stdlib"; instr match { CALL math_Ceil } gen { math.Ceil($x); } bind { x as "x"; } }
pattern math_Round { name: "math.Round"; library: "go-stdlib"; instr match { CALL math_Round } gen { math.Round($x); } bind { x as "x"; } }
pattern math_Log { name: "math.Log"; library: "go-stdlib"; instr match { CALL math_Log } gen { math.Log($x); } bind { x as "x"; } }
pattern math_Exp { name: "math.Exp"; library: "go-stdlib"; instr match { CALL math_Exp } gen { math.Exp($x); } bind { x as "x"; } }
pattern math_Mod { name: "math.Mod"; library: "go-stdlib"; instr match { CALL math_Mod } gen { math.Mod($x, $y); } bind { x as "x"; y as "y"; } }

pattern rand_Intn { name: "math/rand.Intn"; library: "go-stdlib"; instr match { CALL math_rand_Intn } gen { rand.Intn($n); } bind { n as "n"; } }
pattern rand_Float64 { name: "math/rand.Float64"; library: "go-stdlib"; instr match { CALL math_rand_Float64 } gen { rand.Float64(); } }
pattern rand_Int { name: "math/rand.Int"; library: "go-stdlib"; instr match { CALL math_rand_Int } gen { rand.Int(); } }
pattern rand_New { name: "math/rand.New"; library: "go-stdlib"; instr match { CALL math_rand_New } gen { rand.New($src); } bind { src as "src"; } }

pattern sort_Ints { name: "sort.Ints"; library: "go-stdlib"; instr match { CALL sort_Ints } gen { sort.Ints($a); } bind { a as "a"; } }
pattern sort_Strings { name: "sort.Strings"; library: "go-stdlib"; instr match { CALL sort_Strings } gen { sort.Strings($a); } bind { a as "a"; } }
pattern sort_Float64s { name: "sort.Float64s"; library: "go-stdlib"; instr match { CALL sort_Float64s } gen { sort.Float64s($a); } bind { a as "a"; } }
pattern sort_Slice { name: "sort.Slice"; library: "go-stdlib"; instr match { CALL sort_Slice } gen { sort.Slice($s, $less); } bind { s as "s"; less as "less"; } }
pattern sort_Search { name: "sort.Search"; library: "go-stdlib"; instr match { CALL sort_Search } gen { sort.Search($n, $f); } bind { n as "n"; f as "f"; } }
pattern sort_Stable { name: "sort.Stable"; library: "go-stdlib"; instr match { CALL sort_Stable } gen { sort.Stable($data); } bind { data as "data"; } }
pattern sort_IsSorted { name: "sort.IsSorted"; library: "go-stdlib"; instr match { CALL sort_IsSorted } gen { sort.IsSorted($data); } bind { data as "data"; } }
