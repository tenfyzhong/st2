enum EEE {
    A = 1;
    B = 2;
}

// hhhh
struct SS { // aa
    // ss
    1: optional bool a, // jjj
    2: byte b,
    3: i16 c,
    4: i32 d,
    5: i64 e,
    6: double f,
    7: string g,
    8: binary h,
    9: map<i32, i32> i,
    10: optional list<i32> j,
    11: set<i32> k,
}

struct AAA {
    1: string hello,
}

struct BBB {
    1: i16 b1,
    2: i32 b2,
    3: EEE e,
    4: map<AAA, BBB> mapab,
    5: set<AAA> seta,
    6: list<BBB> listb,
}

union UUU {
    1: AAA a;
    2: BBB b;
}
