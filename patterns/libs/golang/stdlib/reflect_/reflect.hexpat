arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_TypeOf { name: "reflect.TypeOf"; library: "go-stdlib"; instr match { CALL reflect_TypeOf } gen { reflect.TypeOf($i); } bind { i as "i"; } }
pattern reflect_ValueOf { name: "reflect.ValueOf"; library: "go-stdlib"; instr match { CALL reflect_ValueOf } gen { reflect.ValueOf($i); } bind { i as "i"; } }
pattern reflect_DeepEqual { name: "reflect.DeepEqual"; library: "go-stdlib"; instr match { CALL reflect_DeepEqual } gen { reflect.DeepEqual($x, $y); } bind { x as "x"; y as "y"; } }
pattern reflect_Value_Interface { name: "reflect.Value.Interface"; library: "go-stdlib"; instr match { CALL reflect___Value__Interface } gen { $v.Interface(); } bind { v as "v"; } }
pattern reflect_Value_Int { name: "reflect.Value.Int"; library: "go-stdlib"; instr match { CALL reflect___Value__Int } gen { $v.Int(); } bind { v as "v"; } }
pattern reflect_Value_String { name: "reflect.Value.String"; library: "go-stdlib"; instr match { CALL reflect___Value__String } gen { $v.String(); } bind { v as "v"; } }
pattern reflect_Value_Bool { name: "reflect.Value.Bool"; library: "go-stdlib"; instr match { CALL reflect___Value__Bool } gen { $v.Bool(); } bind { v as "v"; } }
pattern reflect_Value_Float { name: "reflect.Value.Float"; library: "go-stdlib"; instr match { CALL reflect___Value__Float } gen { $v.Float(); } bind { v as "v"; } }
pattern reflect_Value_Len { name: "reflect.Value.Len"; library: "go-stdlib"; instr match { CALL reflect___Value__Len } gen { $v.Len(); } bind { v as "v"; } }
pattern reflect_Value_Index { name: "reflect.Value.Index"; library: "go-stdlib"; instr match { CALL reflect___Value__Index } gen { $v.Index($i); } bind { v as "v"; i as "i"; } }
pattern reflect_Value_Field { name: "reflect.Value.Field"; library: "go-stdlib"; instr match { CALL reflect___Value__Field } gen { $v.Field($i); } bind { v as "v"; i as "i"; } }
pattern reflect_Value_MapIndex { name: "reflect.Value.MapIndex"; library: "go-stdlib"; instr match { CALL reflect___Value__MapIndex } gen { $v.MapIndex($key); } bind { v as "v"; key as "key"; } }
pattern reflect_Value_MapKeys { name: "reflect.Value.MapKeys"; library: "go-stdlib"; instr match { CALL reflect___Value__MapKeys } gen { $v.MapKeys(); } bind { v as "v"; } }
pattern reflect_Value_Set { name: "reflect.Value.Set"; library: "go-stdlib"; instr match { CALL reflect___Value__Set } gen { $v.Set($x); } bind { v as "v"; x as "x"; } }
pattern reflect_Value_SetString { name: "reflect.Value.SetString"; library: "go-stdlib"; instr match { CALL reflect___Value__SetString } gen { $v.SetString($x); } bind { v as "v"; x as "x"; } }
pattern reflect_Value_Elem { name: "reflect.Value.Elem"; library: "go-stdlib"; instr match { CALL reflect___Value__Elem } gen { $v.Elem(); } bind { v as "v"; } }
pattern reflect_Value_Kind { name: "reflect.Value.Kind"; library: "go-stdlib"; instr match { CALL reflect___Value__Kind } gen { $v.Kind(); } bind { v as "v"; } }
pattern reflect_Value_Slice { name: "reflect.Value.Slice"; library: "go-stdlib"; instr match { CALL reflect___Value__Slice } gen { $v.Slice($i, $j); } bind { v as "v"; i as "i"; j as "j"; } }
pattern reflect_Value_Close { name: "reflect.Value.Close"; library: "go-stdlib"; instr match { CALL reflect___Value__Close } gen { $v.Close(); } bind { v as "v"; } }
pattern reflect_Value_IsNil { name: "reflect.Value.IsNil"; library: "go-stdlib"; instr match { CALL reflect___Value__IsNil } gen { $v.IsNil(); } bind { v as "v"; } }
pattern reflect_Value_IsValid { name: "reflect.Value.IsValid"; library: "go-stdlib"; instr match { CALL reflect___Value__IsValid } gen { $v.IsValid(); } bind { v as "v"; } }
pattern reflect_Type_Name { name: "reflect.Type.Name"; library: "go-stdlib"; instr match { CALL reflect___Type__Name } gen { $t.Name(); } bind { t as "t"; } }
pattern reflect_Type_Kind { name: "reflect.Type.Kind"; library: "go-stdlib"; instr match { CALL reflect___Type__Kind } gen { $t.Kind(); } bind { t as "t"; } }
pattern reflect_Type_NumMethod { name: "reflect.Type.NumMethod"; library: "go-stdlib"; instr match { CALL reflect___Type__NumMethod } gen { $t.NumMethod(); } bind { t as "t"; } }
pattern reflect_Type_NumField { name: "reflect.Type.NumField"; library: "go-stdlib"; instr match { CALL reflect___Type__NumField } gen { $t.NumField(); } bind { t as "t"; } }
pattern reflect_Type_Field { name: "reflect.Type.Field"; library: "go-stdlib"; instr match { CALL reflect___Type__Field } gen { $t.Field($i); } bind { t as "t"; i as "i"; } }
pattern reflect_Type_Size { name: "reflect.Type.Size"; library: "go-stdlib"; instr match { CALL reflect___Type__Size } gen { $t.Size(); } bind { t as "t"; } }
pattern reflect_Type_String { name: "reflect.Type.String"; library: "go-stdlib"; instr match { CALL reflect___Type__String } gen { $t.String(); } bind { t as "t"; } }
pattern reflect_Type_Elem { name: "reflect.Type.Elem"; library: "go-stdlib"; instr match { CALL reflect___Type__Elem } gen { $t.Elem(); } bind { t as "t"; } }
