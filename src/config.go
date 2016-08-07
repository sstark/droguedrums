
package main

type Kit []struct {
    Key string
    Channel int
    Note int
}

type Part struct {
    Name string
    Set string
    Step int
    Bpm int
    Lanes matrix
}
