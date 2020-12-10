package LY

type protocol interface {
	GetName()string
	Input(buffer []byte) int
	Encode(buffer []byte)
	Decode(data string)
}

