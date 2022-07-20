# smt_proof

Test cosmos' SMT (v2alpha1) proof with `ibc-go` verification functions `VerifyMembership` and `VerifyNonMembership`
- Using the `multi` store with SMT store
- Failed with the current ics23 package
    - The `ibc-go` verification requires the original(non-hashed) key and value, however SMT proof has the hashed key
    - The [fix](https://github.com/confio/ics23/pull/88) has been proposed
