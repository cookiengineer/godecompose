arch x86_64;
platform linux, darwin, windows, freebsd;

pattern crypto_sha256_Sum256 {
    name: "crypto/sha256.Sum256";
    library: "go-stdlib";
    description: "Computes SHA256 hash";

    instr match {
        MOVQ data, AX
        CALL sha256_Sum256
    }

    gen {
        sha256.Sum256($data);
    }

    bind {
        data as "data";
    }
}
