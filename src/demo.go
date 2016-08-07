package main

func demo1() Part {
    l := Part{
        Name: "demo1",
        Bpm: 160,
        Step: 4,
        Lanes: matrix{
            {1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
            {6, 0, 6, 6, 8, 0, 6, 0, 6, 6, 6, 0},
            {0, 0, 12, 0, 0, 2, 5, 0, 3, 0, 0, 2},
        },
    }
    return l
}

func demo2() Part {
    l := Part{
        Name: "demo2",
        Bpm: 160,
        Step: 8,
        Lanes: matrix{
            {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
            {6, 0, 6, 6, 8, 0, 6, 0, 6, 6, 6, 0},
            {0, 12, 0, 9, 12, 2, 5, 12, 3, 9, 12, 2},
        },
    }
    return l
}
