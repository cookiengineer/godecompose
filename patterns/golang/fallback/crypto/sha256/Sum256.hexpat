arch x86_64;
platform linux, darwin, windows, freebsd;

pattern crypto_sha256_Sum256_call {
    name: "crypto/sha256.Sum256";
    library: "go-stdlib";
    description: "sha256.Sum256 call";

    instr match {
        CALL sha256_Sum256
    }

    gen {
        h := sha256.Sum256($data);
    }

    bind {
        data as "data";
    }
}
