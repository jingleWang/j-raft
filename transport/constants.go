package transport

import "time"

const MagicVal uint16 = 2532
const Version byte = 1
const HeaderLength uint64 = 16

const DefaultTimeOut = 3 * time.Second
