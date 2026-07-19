arch x86_64;
platform linux, darwin, windows, freebsd;

pattern crypto_sha256_New {
    name: "crypto/sha256.New";
    library: "go-stdlib";
    description: "Creates a new SHA256 hash";

    instr match {
        CALL sha256_New
    }

    gen {
        sha256.New();
    }
}
