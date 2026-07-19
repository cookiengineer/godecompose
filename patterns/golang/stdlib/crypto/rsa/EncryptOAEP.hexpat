arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rsa_EncryptOAEP { name: "crypto/rsa.EncryptOAEP"; library: "go-stdlib"; instr match { CALL crypto_rsa_EncryptOAEP } gen { rsa.EncryptOAEP($hash, $random, $pub, $msg, $label); } bind { hash as "hash"; random as "random"; pub as "pub"; msg as "msg"; label as "label"; } }
