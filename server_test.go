package main

import (
    "testing"
)

func add(a, b int) int {
    return a + b
}

func TestAdd(t *testing.T) {
    type args struct {
        a int
        b int
    }

    tests := []struct {
        name string
        args args
        want int
    }{
        {
            name: "normal",
            args: args{a: 1, b: 2},
            want: 3,
        },
        {
            name: "another",
            args: args{a: 100, b: 1},
            want: 101,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            actual := add(tt.args.a, tt.args.b)

            if actual != tt.want {
                t.Errorf("add() = %v, want %v", actual, tt.want)
            }
        })
    }
}
