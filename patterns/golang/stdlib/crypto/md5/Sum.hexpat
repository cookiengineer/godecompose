arch x86_64;
platform linux, darwin, windows, freebsd;

pattern crypto_md5_Sum {
    name: "crypto/md5.Sum";
    library: "go-stdlib";
    description: "Computes MD5 hash";

    instr match {
        MOVQ data, AX
        CALL md5_Sum
    }

    gen {
        md5.Sum($data);
    }

    bind {
        data as "data";
    }
}
