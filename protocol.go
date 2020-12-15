package ly

type protocol interface {
	GetName()string  // 获取协议名称
	Input(buffer []byte) int // 获取一条信息长度
	Decode(buffer []byte) interface{}  // 解码信息
	Encode(data string)string     // 加密信息
}

