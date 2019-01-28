package x

// OpenRedPacket opens a red packet,
// amount and quantity are remaining data.
// amount and quantity must > 0.
// amount must >= quantity.
func OpenRedPacket(amount, quantity int64) int64 {
	if quantity == 1 {
		return amount
	}
	i := amount / quantity
	n := Random(i, i*2)
	return n
}
