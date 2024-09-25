package cbor

var tests_ExampleEncoded = []struct {
	encoded string
}{
	{encoded: "00"},
	{encoded: "01"},
	{encoded: "0a"},
	{encoded: "17"},
	{encoded: "1818"},
	{encoded: "1819"},
	{encoded: "1864"},
	{encoded: "1903e8"},
	{encoded: "1a000f4240"},
	{encoded: "1b000000e8d4a51000"},
	{encoded: "1bffffffffffffffff"},
	{encoded: "c249010000000000000000"},
	{encoded: "3bffffffffffffffff"},
	{encoded: "c349010000000000000000"},
	{encoded: "20"},
	{encoded: "29"},
	{encoded: "3863"},
	{encoded: "3903e7"},
	{encoded: "f90000"},
	{encoded: "f98000"},
	{encoded: "f93c00"},
	{encoded: "fb3ff199999999999a"},
	{encoded: "f93e00"},
	{encoded: "f97bff"},
	{encoded: "fa47c35000"},
	{encoded: "fa7f7fffff"},
	{encoded: "fb7e37e43c8800759c"},
	{encoded: "f90001"},
	{encoded: "f90400"},
	{encoded: "f9c400"},
	{encoded: "fbc010666666666666"},
	{encoded: "f97c00"},
	{encoded: "f97e00"},
	{encoded: "f9fc00"},
	{encoded: "fa7f800000"},
	{encoded: "fa7fc00000"},
	{encoded: "faff800000"},
	{encoded: "fb7ff0000000000000"},
	{encoded: "fb7ff8000000000000"},
	{encoded: "fbfff0000000000000"},
	{encoded: "f4"},
	{encoded: "f5"},
	{encoded: "f6"},
	{encoded: "f7"},
	{encoded: "f0"},
	{encoded: "f8ff"},
	{encoded: "c074323031332d30332d32315432303a30343a30305a"},
	{encoded: "c11a514b67b0"},
	{encoded: "c1fb41d452d9ec200000"},
	{encoded: "d74401020304"},
	{encoded: "d818456449455446"},
	{encoded: "d82076687474703a2f2f7777772e6578616d706c652e636f6d"},
	{encoded: "40"},
	{encoded: "4401020304"},
	{encoded: "60"},
	{encoded: "6161"},
	{encoded: "6449455446"},
	{encoded: "62225c"},
	{encoded: "62c3bc"},
	{encoded: "63e6b0b4"},
	{encoded: "64f0908591"},
	{encoded: "80"},
	{encoded: "83010203"},
	{encoded: "8301820203820405"},
	{encoded: "98190102030405060708090a0b0c0d0e0f101112131415161718181819"},
	{encoded: "a0"},
	{encoded: "a201020304"},
	{encoded: "a26161016162820203"},
	{encoded: "826161a161626163"},
	{encoded: "a56161614161626142616361436164614461656145"},
	{encoded: "5f42010243030405ff"},
	{encoded: "7f657374726561646d696e67ff"},
	{encoded: "9fff"},
	{encoded: "9f018202039f0405ffff"},
	{encoded: "9f01820203820405ff"},
	{encoded: "83018202039f0405ff"},
	{encoded: "83019f0203ff820405"},
	{encoded: "9f0102030405060708090a0b0c0d0e0f101112131415161718181819ff"},
	{encoded: "bf61610161629f0203ffff"},
	{encoded: "826161bf61626163ff"},
	{encoded: "bf6346756ef563416d7421ff"},
}
