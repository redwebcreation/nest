package docker

//func nextIP(ip net.IP, increment uint) net.IP {
//	i := ip.To4()
//	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
//	v += increment
//	v3 := byte(v & 0xFF)
//	v2 := byte((v >> 8) & 0xFF)
//	v1 := byte((v >> 16) & 0xFF)
//	v0 := byte((v >> 24) & 0xFF)
//	return net.IPv4(v0, v1, v2, v3)
//}
