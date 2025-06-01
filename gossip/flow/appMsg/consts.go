package appMsg

const (
	HandshakeMsg = iota
	//Not done
	SendInvBlockHash
	SendProposeBlock

	RequestHighestBlockHash
	ResponseHighestBlockHash

	RequestInvBlock
	ResponseInvBlock

	RequestBlock
	ResponseBlock

	RequestBlockHashRelay
	ResponseBlockHashRelay

	RequestTransaction
	ResponseTransaction
)
