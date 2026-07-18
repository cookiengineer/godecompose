arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tar_NewReader { name: "archive/tar.NewReader"; library: "go-stdlib"; instr match { CALL archive_tar_NewReader } gen { tar.NewReader($r); } bind { r as "r"; } }
pattern tar_NewWriter { name: "archive/tar.NewWriter"; library: "go-stdlib"; instr match { CALL archive_tar_NewWriter } gen { tar.NewWriter($w); } bind { w as "w"; } }
pattern tar_Reader_Next { name: "tar.Reader.Next"; library: "go-stdlib"; instr match { CALL archive_tar___Reader__Next } gen { $r.Next(); } bind { r as "r"; } }
pattern tar_Reader_Read { name: "tar.Reader.Read"; library: "go-stdlib"; instr match { CALL archive_tar___Reader__Read } gen { $r.Read($b); } bind { r as "r"; b as "b"; } }
pattern tar_Writer_WriteHeader { name: "tar.Writer.WriteHeader"; library: "go-stdlib"; instr match { CALL archive_tar___Writer__WriteHeader } gen { $w.WriteHeader($hdr); } bind { w as "w"; hdr as "hdr"; } }
pattern tar_Writer_Write { name: "tar.Writer.Write"; library: "go-stdlib"; instr match { CALL archive_tar___Writer__Write } gen { $w.Write($b); } bind { w as "w"; b as "b"; } }
pattern tar_Writer_Close { name: "tar.Writer.Close"; library: "go-stdlib"; instr match { CALL archive_tar___Writer__Close } gen { $w.Close(); } bind { w as "w"; } }
pattern tar_Writer_Flush { name: "tar.Writer.Flush"; library: "go-stdlib"; instr match { CALL archive_tar___Writer__Flush } gen { $w.Flush(); } bind { w as "w"; } }

pattern zip_OpenReader { name: "archive/zip.OpenReader"; library: "go-stdlib"; instr match { CALL archive_zip_OpenReader } gen { zip.OpenReader($name); } bind { name as "name"; } }
pattern zip_NewReader { name: "archive/zip.NewReader"; library: "go-stdlib"; instr match { CALL archive_zip_NewReader } gen { zip.NewReader($r, $size); } bind { r as "r"; size as "size"; } }
pattern zip_NewWriter { name: "archive/zip.NewWriter"; library: "go-stdlib"; instr match { CALL archive_zip_NewWriter } gen { zip.NewWriter($w); } bind { w as "w"; } }
pattern zip_File_Open { name: "zip.File.Open"; library: "go-stdlib"; instr match { CALL archive_zip___File__Open } gen { $f.Open(); } bind { f as "f"; } }
pattern zip_Writer_Create { name: "zip.Writer.Create"; library: "go-stdlib"; instr match { CALL archive_zip___Writer__Create } gen { $w.Create($name); } bind { w as "w"; name as "name"; } }
pattern zip_Writer_CreateHeader { name: "zip.Writer.CreateHeader"; library: "go-stdlib"; instr match { CALL archive_zip___Writer__CreateHeader } gen { $w.CreateHeader($fh); } bind { w as "w"; fh as "fh"; } }
pattern zip_Writer_Close { name: "zip.Writer.Close"; library: "go-stdlib"; instr match { CALL archive_zip___Writer__Close } gen { $w.Close(); } bind { w as "w"; } }
