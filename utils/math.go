package utils

func Min(x, y int) int {
    if x < y {
        return x
    }
    return y
}

func Max(x, y int) int {
    if x > y {
        return x
    }
    return y
}

func Abs(n int) int {
    if n < 0 {
       n = -n
    }
    return n
}
