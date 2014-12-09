package main

import "fmt"
import "encoding/hex"
import "encoding/base64"
import "math"
import "io/ioutil"
import "strings"

func hexToBase64(in string) (string, error) {
	bytes, err :=  hex.DecodeString(in)
	if err != nil {
		return "", err
	}

	out := base64.StdEncoding.EncodeToString(bytes)
	return out, nil
}

func sameLengthXor(b1 []byte, b2 []byte) ([]byte) {
	b3 := make([]byte,len(b1))
	for i := 0; i < len(b1); i++ {
		b3[i] = b1[i] ^ b2[i]
	}
	return b3
}

func repeatKeyXor(payload []byte, key []byte) ([]byte) {
	out := make([]byte,len(payload));
	payLength := len(payload)
	keyLength := len(key)
	for i := 0; i < payLength; i++ {
		keyI := i % keyLength
		out[i] = payload[i] ^ key[keyI]
	}
	return out
}

func singleByteXor(b1 []byte, b2 byte) ([]byte) {
	b3 := make([]byte,len(b1))
	for i := 0; i < len(b1); i ++ {
		b3[i] = b1[i] ^ b2
	}
	return b3
}

func distributionDiff(text string) (float64) {
	englishLetterStats := map[rune]float64 {
		'e': 12.702, 't': 9.056, 'a': 8.167,
		'o': 7.507, 'i': 6.996, 'n': 6.749,
		's': 6.237, 'h': 6.094, 'r': 5.987,
		'd': 4.253, 'l': 4.025, 'u': 2.578,
	}

	counts := map[rune]float64 {
		'e': 0, 't': 0, 'a': 0, 'o': 0, 'i': 0, 'n': 0,
		's': 0, 'h': 0, 'r': 0, 'd': 0, 'l': 0, 'u': 0,
	}

	totalCounted := 0
	diff := 0.0

	for _,c := range(text) {
		_, ok := englishLetterStats[c]
		if ok {
			counts[c] = counts[c] + 1
			totalCounted += 1
		}
	}

	if totalCounted == 0 {
		return math.MaxInt32
	}

	for c, freq := range(englishLetterStats) {
		diff += math.Pow((freq - counts[c]/float64(totalCounted) * 100),2)
	}

	return diff
}

func bestSingleCharXor(cipherBytes []byte) (byte, float64, string) {
	minI, minScore, minOut := byte(0), 1000000.0, ""

	for i := byte(33); i < 133; i++ {
		out := string(singleByteXor(cipherBytes, i))
		score := distributionDiff(out)
		if score < math.MaxInt32 {
			//fmt.Printf("%c %f\n",i,score)
			if score < minScore {
				minI = i
				minScore = score
				minOut = out
			}
		}
	}


	return minI, minScore, minOut
}

func bitDiff(b1 byte, b2 byte) (int) {
	equal := 0
	for i := uint(0); i < 8; i++ {
		b1_bit := (b1 & (1 << i)) >> i
		b2_bit := (b2 & (1 << i)) >> i
		if b1_bit == b2_bit {
			equal += 1
		}
	}
	return (8-equal)
}

func editDistance(s1 []byte, s2 []byte) (int) {
	distance := 0
	for n, b1 := range(s1) {
		distance += bitDiff(b1, s2[n])
	}
	return distance
}

func main() {
	//fmt.Println(hexToBase64("49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"))

	//b1,_ := hex.DecodeString("1c0111001f010100061a024b53535009181c")
	//b2,_ := hex.DecodeString("686974207468652062756c6c277320657965")
	//b3 := sameLengthXor(b1,b2)

//	cipherBytes,_ := hex.DecodeString("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")

	minI, minScore, minOut := byte(0), 1000000.0, ""
		//bestSingleCharXor(cipherBytes)
	minLineNo := -1

	fileBytes, _ := ioutil.ReadFile("4.txt")
	lines := strings.Split(string(fileBytes),"\n")
	for lineNo,line := range(lines) {
		lineBytes, _ := hex.DecodeString(line)
		i,score,out := bestSingleCharXor(lineBytes)
		if score < minScore {
			minI = i
			minScore = score
			minOut = out
			minLineNo = lineNo
		}
	}

	fmt.Printf("score = %f, minI = %c, minOut = %s, lineNo = %d\n",minScore,minI,minOut,minLineNo)

	str := []byte(`Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`)
	key := []byte("ICE")
	xorBytes := repeatKeyXor(str,key)
	xorStr := hex.EncodeToString(xorBytes)

	fmt.Printf("repeat key xor = %s\n", xorStr)
	// fmt.Printf("bitDiffs (255,0)=%d\n", bitDiff(255,0))
	// fmt.Printf("bitDiffs (127,255)=%d\n", bitDiff(127,255))
	// fmt.Printf("bitDiffs (16,4)=%d\n", bitDiff(16,4))
	fmt.Printf("editDistance ('this is a test','wokka wokka!!!')=%d\n", editDistance([]byte("this is a test"),[]byte("wokka wokka!!!")))
}
