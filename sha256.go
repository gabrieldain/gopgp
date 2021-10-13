package main

import (
    "fmt"
    "strings"
    "strconv"
)

const (
    h0 = 0x6A09E667
    h1 = 0xBB67AE85
    h2 = 0x3C6EF372
    h3 = 0xA54FF53A
    h4 = 0x510E527F
    h5 = 0x9B05688C
    h6 = 0x1F83D9AB
    h7 = 0x5BE0CD19
)

var _K = []uint32{
    0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1,
    0x923f82a4, 0xab1c5ed5, 0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3,
    0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174, 0xe49b69c1, 0xefbe4786,
    0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
    0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147,
    0x06ca6351, 0x14292967, 0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13,
    0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85, 0xa2bfe8a1, 0xa81a664b,
    0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
    0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a,
    0x5b9cca4f, 0x682e6ff3, 0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208,
    0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
}

func and(s1 string, s2 string) (ans string) {
    for i := 0; i < len(s1); i++ {
        if s1[i] == '1' && s2[i] == '1' {
            ans = ans + "1"
        } else {
            ans = ans + "0"
        }
    }
    return
}

func not(s string) (ans string) {
    for i:= range s {
        if s[i] == '0' {
            ans = ans + "1"
        } else {
            ans = ans + "0"
        }
    }
    return
}


func xor(s1 string, s2 string) (ans string) {
    //x1, _ := strconv.ParseInt(s1, 2, 64)
    //x2, _ := strconv.ParseInt(s2, 2, 64)
    //ans = fmt.Sprintf("%.32b", x1^x2)
    //return
    for i := 0; i < len(s1); i++ {
        if s1[i] == s2[i] {
            ans = ans + "0"
        } else {
            ans = ans + "1"
        }
    }
    return
}

func rotateRight(s string, l int) string {
    for l > 0 {
        s = s[len(s)-1:] + s[:len(s)-1]
        l = l - 1
    }
    return s
}

func rotateLeft (s string, l int) string {
    for l > 0 {
        s = s[1:] + s[:1]
        l = l - 1
    }
    return s
}

func shiftRight(s string, l int) string {
    x, _ := strconv.ParseInt(s, 2, 64)
    x = x>>l
    return fmt.Sprintf("%.32b", x)
}

func shiftLeft (s string, l int) string {
    x, _ := strconv.ParseInt(s, 2, 64)
    x = x<<l
    return fmt.Sprintf("%.32b", x)
}

// isBitString checks if the input is a binary number encoded as a string.
func isBitString(s string) bool {
    s = strings.ReplaceAll(s, "1", "")
    s = strings.ReplaceAll(s, "0", "")
    if s == "" {
        return true
    } else {
        return false
    }
}

// stringToBits returns the binary representation of the input as a string.
func stringToBits(s string) (bitString string) {
    for _, char := range s {
        bitString = fmt.Sprintf("%s%.8b", bitString, char)
    }
    return
}

// bitsToInt returns the decimal representation of a bitstring as an uint64.
func bitsToInt(s string) uint64 {
    x, _ := strconv.ParseUint(s, 2, 64)
    return x
}

// preProcess adds a "1" to the end of the input, followed by a number of
// zeroes such that the length of the original string plus the added zeroes
// (not including the added 1s) is a multiple of 512, minus 64, followed by
// the length of the original input encoded as a 64-bit binary.
func preProcess(s string) string {
    if isBitString(s) == false {
        s = stringToBits(s)
    }
    length := fmt.Sprintf("%.64b", len(s))
    padding := strings.Repeat("0", 512 - ((len(s) + 1 + 64) % 512))
    s = s+"1"+padding+length
    return s
}

func createChunks(s string) [][]string {
    var (
        chunks [][]string
        words []string
        currentLen int
        currentStart int
    )
    for i := range s {
        if currentLen == 512 {
            chunks = append(chunks, []string{s[currentStart:i]})
            currentLen = 0
            currentStart = i
        }
        currentLen++
    }
    chunks = append(chunks, []string{s[currentStart:]})
    for i := range chunks {
        currentLen = 0
        currentStart = 0
        chunks[i][0] = chunks[i][0] + strings.Repeat("0", 1536)
        words = []string{}
        for j := range chunks[i][0] {
            if currentLen == 32 {
                words = append(words, chunks[i][0][currentStart:j])
                currentLen = 0
                currentStart = j
            }
            currentLen++
        }
        words = append(words, chunks[i][0][currentStart:])
        chunks[i] = words
    }
    return chunks
}

func createMessageSchedule(chunks [][]string) [][]string {
    var (
        op1 string
        op2 string
        op3 string
        s0 string
        s1 string
        s0_int uint64
        s1_int uint64
        sum uint64
        x uint64
        y uint64
    )
    for i := range chunks {
        for j := 16; j < 64; j++{
            op1 = rotateRight(chunks[i][j-15], 7)
            op2 = rotateRight(chunks[i][j-15], 18)
            op3 = shiftRight(chunks[i][j-15], 3)
            s0 = xor(xor(op1, op2), op3)
            s0_int, _ = strconv.ParseUint(s0, 2, 64)
            op1 = rotateRight(chunks[i][j-2], 17)
            op2 = rotateRight(chunks[i][j-2], 19)
            op3 = shiftRight(chunks[i][j-2], 10)
            s1 = xor(xor(op1, op2), op3)
            s1_int, _ = strconv.ParseUint(s1, 2, 64)
            x, _ = strconv.ParseUint(chunks[i][j-16], 2, 64)
            y, _ = strconv.ParseUint(chunks[i][j-7], 2, 64)
            sum = s0_int + s1_int + x + y
            sum = sum % (1<<32) // sum is calculated modulo 2^32
            chunks[i][j] = fmt.Sprintf("%.32b", sum)
        }
    }
    return chunks
}

func compressMessage(chunks [][]string) string {
    var (
        h_0, h_1, h_2, h_3, h_4, h_5, h_6, h_7 int = h0,h1,h2,h3,h4,h5,h6,h7
        s0, s1, op1, op2, op3, ch, maj string
        s0_int, s1_int, d_int, h_int, w_int, ch_int, maj_int, tmp1, tmp2 uint64
    )
    for i  := range chunks {
        a := fmt.Sprintf("%.32b", h_0)
        b := fmt.Sprintf("%.32b", h_1)
        c := fmt.Sprintf("%.32b", h_2)
        d := fmt.Sprintf("%.32b", h_3)
        e := fmt.Sprintf("%.32b", h_4)
        f := fmt.Sprintf("%.32b", h_5)
        g := fmt.Sprintf("%.32b", h_6)
        h := fmt.Sprintf("%.32b", h_7)
        for j := range chunks[i] {
            op1 = rotateRight(e, 6)
            op2 = rotateRight(e, 11)
            op3 = rotateRight(e, 25)
            s1 = xor(xor(op1, op2), op3)
            ch = xor(and(e, f), and(not(e), g))
            h_int, _ = strconv.ParseUint(h, 2, 64)
            s1_int, _ = strconv.ParseUint(s1, 2, 64)
            ch_int, _ = strconv.ParseUint(ch, 2, 64)
            w_int, _ = strconv.ParseUint(chunks[i][j], 2, 64)
            tmp1 = h_int + s1_int + ch_int + uint64(_K[j]) + w_int
            tmp1 = tmp1 % (1<<32) // tmp1 is calculated modulo 2^32
            op1 = rotateRight(a, 2)
            op2 = rotateRight(a, 13)
            op3 = rotateRight(a, 22)
            s0 = xor(xor(op1, op2), op3)
            maj = xor(xor(and(a, b), and(a, c)), and(b, c))
            d_int, _ = strconv.ParseUint(d, 2, 64)
            s0_int, _ = strconv.ParseUint(s0, 2, 64)
            maj_int, _ = strconv.ParseUint(maj, 2, 64)
            tmp2 = s0_int + maj_int
            tmp2 = tmp2 % (1<<32) // tmp2 is calculated modulo 2^32
            h = g
            g = f
            f = e
            e = fmt.Sprintf("%.32b", (d_int + tmp1) % (1<<32))
            d = c
            c = b
            b = a
            a = fmt.Sprintf("%.32b", (tmp1 + tmp2) % (1<<32))
        }
        h_0 = (h_0 + int(bitsToInt(a))) % (1<<32)
        h_1 = (h_1 + int(bitsToInt(b))) % (1<<32)
        h_2 = (h_2 + int(bitsToInt(c))) % (1<<32)
        h_3 = (h_3 + int(bitsToInt(d))) % (1<<32)
        h_4 = (h_4 + int(bitsToInt(e))) % (1<<32)
        h_5 = (h_5 + int(bitsToInt(f))) % (1<<32)
        h_6 = (h_6 + int(bitsToInt(g))) % (1<<32)
        h_7 = (h_7 + int(bitsToInt(h))) % (1<<32)
    }
    h_0_hex := fmt.Sprintf("%x", h_0)
    h_1_hex := fmt.Sprintf("%x", h_1)
    h_2_hex := fmt.Sprintf("%x", h_2)
    h_3_hex := fmt.Sprintf("%x", h_3)
    h_4_hex := fmt.Sprintf("%x", h_4)
    h_5_hex := fmt.Sprintf("%x", h_5)
    h_6_hex := fmt.Sprintf("%x", h_6)
    h_7_hex := fmt.Sprintf("%x", h_7)
    return h_0_hex+h_1_hex+h_2_hex+h_3_hex+h_4_hex+h_5_hex+h_6_hex+h_7_hex
}

func Digest(s string) string {
    return compressMessage(createMessageSchedule(createChunks(preProcess(s))))
}

func main() {
    var s string
    s = "hello world"
    //s = strings.Repeat("A", 10000)
    fmt.Println(Digest(s))
}
