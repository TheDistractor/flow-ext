//Package conversions provides Gadgets that offer some form of data conversion process
package conversions

import (
	"fmt"
	"testing"
	"time"
	"github.com/golang/glog"
	"github.com/jcw/flow"

)

func init() {
	glog.Infoln("mapping tests...")
}





type MockOutput int

func (c *MockOutput) Send(m flow.Message) {
	fmt.Printf("%T: %v\n", m, m)
}

func (c *MockOutput) Disconnect() {}


func ExampleRangeFlow() {

	g := new(RangeMapper)
	p := make(chan flow.Message, 5)
	in := make(chan flow.Message, 1)
	g.Param, g.In, g.Out = p, in, new(MockOutput)

	p <- flow.Tag{"fromlow", float64(0)}
	p <- flow.Tag{"fromhi", float64(255)}
	p <- flow.Tag{"tolow", float64(0)}
	p <- flow.Tag{"tohi", float64(1023)}



	go func() { g.Run() }()
	close(p)


	in <- float64(127)
	in <- float64(128)

	in <- flow.Tag{"hello", float64(10)}

	in <- flow.Tag{"world", float64(100)}

	in <- flow.Tag{"hi", float64(254)}

	in <- flow.Tag{"again", int16(2000)}


	<-time.After(time.Second*1)

// Output:
//float64: 509
//float64: 514
//flow.Tag: {hello 40}
//flow.Tag: {world 401}
//flow.Tag: {hi 1019}
//flow.Tag: {again 1023}



}



func X() {
	fmt.Println(time.Now())
	fmt.Println( flow.Tag{"a",1}  )
}


func TestRangeInvalid(t *testing.T) () {

	r := NewRangeMap()

	if r.Valid() {
		t.Error("NewRangeMap should be Invalid")
	}

}

func TestRangeValid(t *testing.T) () {

	r := NewRangeMap()

	r.SetFromLow(0)
	r.SetFromHi(255)

	r.SetToLow(0)
	r.SetToHi(1023)

	if !r.Valid() {
		t.Error("NewRangeMap should be Valid")
	}

}

func TestRangePartInvalidFrom(t *testing.T) () {

	r := NewRangeMap()

	r.SetFromLow(0)
	r.SetFromHi(255)

	if r.Valid() {
		t.Error("NewRangeMap should not be Valid with partial data")
	}
}

func TestRangePartInvalidTo(t *testing.T) () {

	r := NewRangeMap()

	r.SetToLow(0)
	r.SetToHi(1023)

	if r.Valid() {
		t.Error("NewRangeMap should not be Valid with partial data")
	}
}

func TestNewRangeMapFrom(t *testing.T) () {

	r := NewRangeMapFrom(0, 255, 0, 1023)

	if !r.Valid() {
		t.Error("NewRangeMapFrom should be Valid")
	}
}


func TestNewRangeMapFromCalc(t *testing.T) () {

	r := NewRangeMapFrom( 0,255,0,1023)

	var v int64
	var ok bool

	if v,ok = r.Map(0); !ok {
		t.Error("Map failed")

	}
	if v != 0 {
		t.Error("Expected 0, got", v)
	}


}


func ExampleRangePositive255() {

	a1 := int64(0)
	a2 := int64(255)
	b1 := int64(0)
	b2 := int64(1023)

	for i:=int64(0);i<256;i++ {
		r := NewRangeMapFrom(a1, a2, b1, b2)
		fmt.Printf("\tInput %d Output %d\n",i, r.MustMap(i))
	}
	//  Output:
	//	Input 0 Output 0
	//	Input 1 Output 4
	//	Input 2 Output 8
	//	Input 3 Output 12
	//	Input 4 Output 16
	//	Input 5 Output 20
	//	Input 6 Output 24
	//	Input 7 Output 28
	//	Input 8 Output 32
	//	Input 9 Output 36
	//	Input 10 Output 40
	//	Input 11 Output 44
	//	Input 12 Output 48
	//	Input 13 Output 52
	//	Input 14 Output 56
	//	Input 15 Output 60
	//	Input 16 Output 64
	//	Input 17 Output 68
	//	Input 18 Output 72
	//	Input 19 Output 76
	//	Input 20 Output 80
	//	Input 21 Output 84
	//	Input 22 Output 88
	//	Input 23 Output 92
	//	Input 24 Output 96
	//	Input 25 Output 100
	//	Input 26 Output 104
	//	Input 27 Output 108
	//	Input 28 Output 112
	//	Input 29 Output 116
	//	Input 30 Output 120
	//	Input 31 Output 124
	//	Input 32 Output 128
	//	Input 33 Output 132
	//	Input 34 Output 136
	//	Input 35 Output 140
	//	Input 36 Output 144
	//	Input 37 Output 148
	//	Input 38 Output 152
	//	Input 39 Output 156
	//	Input 40 Output 160
	//	Input 41 Output 164
	//	Input 42 Output 168
	//	Input 43 Output 173
	//	Input 44 Output 177
	//	Input 45 Output 181
	//	Input 46 Output 185
	//	Input 47 Output 189
	//	Input 48 Output 193
	//	Input 49 Output 197
	//	Input 50 Output 201
	//	Input 51 Output 205
	//	Input 52 Output 209
	//	Input 53 Output 213
	//	Input 54 Output 217
	//	Input 55 Output 221
	//	Input 56 Output 225
	//	Input 57 Output 229
	//	Input 58 Output 233
	//	Input 59 Output 237
	//	Input 60 Output 241
	//	Input 61 Output 245
	//	Input 62 Output 249
	//	Input 63 Output 253
	//	Input 64 Output 257
	//	Input 65 Output 261
	//	Input 66 Output 265
	//	Input 67 Output 269
	//	Input 68 Output 273
	//	Input 69 Output 277
	//	Input 70 Output 281
	//	Input 71 Output 285
	//	Input 72 Output 289
	//	Input 73 Output 293
	//	Input 74 Output 297
	//	Input 75 Output 301
	//	Input 76 Output 305
	//	Input 77 Output 309
	//	Input 78 Output 313
	//	Input 79 Output 317
	//	Input 80 Output 321
	//	Input 81 Output 325
	//	Input 82 Output 329
	//	Input 83 Output 333
	//	Input 84 Output 337
	//	Input 85 Output 341
	//	Input 86 Output 345
	//	Input 87 Output 349
	//	Input 88 Output 353
	//	Input 89 Output 357
	//	Input 90 Output 361
	//	Input 91 Output 365
	//	Input 92 Output 369
	//	Input 93 Output 373
	//	Input 94 Output 377
	//	Input 95 Output 381
	//	Input 96 Output 385
	//	Input 97 Output 389
	//	Input 98 Output 393
	//	Input 99 Output 397
	//	Input 100 Output 401
	//	Input 101 Output 405
	//	Input 102 Output 409
	//	Input 103 Output 413
	//	Input 104 Output 417
	//	Input 105 Output 421
	//	Input 106 Output 425
	//	Input 107 Output 429
	//	Input 108 Output 433
	//	Input 109 Output 437
	//	Input 110 Output 441
	//	Input 111 Output 445
	//	Input 112 Output 449
	//	Input 113 Output 453
	//	Input 114 Output 457
	//	Input 115 Output 461
	//	Input 116 Output 465
	//	Input 117 Output 469
	//	Input 118 Output 473
	//	Input 119 Output 477
	//	Input 120 Output 481
	//	Input 121 Output 485
	//	Input 122 Output 489
	//	Input 123 Output 493
	//	Input 124 Output 497
	//	Input 125 Output 501
	//	Input 126 Output 505
	//	Input 127 Output 509
	//	Input 128 Output 514
	//	Input 129 Output 518
	//	Input 130 Output 522
	//	Input 131 Output 526
	//	Input 132 Output 530
	//	Input 133 Output 534
	//	Input 134 Output 538
	//	Input 135 Output 542
	//	Input 136 Output 546
	//	Input 137 Output 550
	//	Input 138 Output 554
	//	Input 139 Output 558
	//	Input 140 Output 562
	//	Input 141 Output 566
	//	Input 142 Output 570
	//	Input 143 Output 574
	//	Input 144 Output 578
	//	Input 145 Output 582
	//	Input 146 Output 586
	//	Input 147 Output 590
	//	Input 148 Output 594
	//	Input 149 Output 598
	//	Input 150 Output 602
	//	Input 151 Output 606
	//	Input 152 Output 610
	//	Input 153 Output 614
	//	Input 154 Output 618
	//	Input 155 Output 622
	//	Input 156 Output 626
	//	Input 157 Output 630
	//	Input 158 Output 634
	//	Input 159 Output 638
	//	Input 160 Output 642
	//	Input 161 Output 646
	//	Input 162 Output 650
	//	Input 163 Output 654
	//	Input 164 Output 658
	//	Input 165 Output 662
	//	Input 166 Output 666
	//	Input 167 Output 670
	//	Input 168 Output 674
	//	Input 169 Output 678
	//	Input 170 Output 682
	//	Input 171 Output 686
	//	Input 172 Output 690
	//	Input 173 Output 694
	//	Input 174 Output 698
	//	Input 175 Output 702
	//	Input 176 Output 706
	//	Input 177 Output 710
	//	Input 178 Output 714
	//	Input 179 Output 718
	//	Input 180 Output 722
	//	Input 181 Output 726
	//	Input 182 Output 730
	//	Input 183 Output 734
	//	Input 184 Output 738
	//	Input 185 Output 742
	//	Input 186 Output 746
	//	Input 187 Output 750
	//	Input 188 Output 754
	//	Input 189 Output 758
	//	Input 190 Output 762
	//	Input 191 Output 766
	//	Input 192 Output 770
	//	Input 193 Output 774
	//	Input 194 Output 778
	//	Input 195 Output 782
	//	Input 196 Output 786
	//	Input 197 Output 790
	//	Input 198 Output 794
	//	Input 199 Output 798
	//	Input 200 Output 802
	//	Input 201 Output 806
	//	Input 202 Output 810
	//	Input 203 Output 814
	//	Input 204 Output 818
	//	Input 205 Output 822
	//	Input 206 Output 826
	//	Input 207 Output 830
	//	Input 208 Output 834
	//	Input 209 Output 838
	//	Input 210 Output 842
	//	Input 211 Output 846
	//	Input 212 Output 850
	//	Input 213 Output 855
	//	Input 214 Output 859
	//	Input 215 Output 863
	//	Input 216 Output 867
	//	Input 217 Output 871
	//	Input 218 Output 875
	//	Input 219 Output 879
	//	Input 220 Output 883
	//	Input 221 Output 887
	//	Input 222 Output 891
	//	Input 223 Output 895
	//	Input 224 Output 899
	//	Input 225 Output 903
	//	Input 226 Output 907
	//	Input 227 Output 911
	//	Input 228 Output 915
	//	Input 229 Output 919
	//	Input 230 Output 923
	//	Input 231 Output 927
	//	Input 232 Output 931
	//	Input 233 Output 935
	//	Input 234 Output 939
	//	Input 235 Output 943
	//	Input 236 Output 947
	//	Input 237 Output 951
	//	Input 238 Output 955
	//	Input 239 Output 959
	//	Input 240 Output 963
	//	Input 241 Output 967
	//	Input 242 Output 971
	//	Input 243 Output 975
	//	Input 244 Output 979
	//	Input 245 Output 983
	//	Input 246 Output 987
	//	Input 247 Output 991
	//	Input 248 Output 995
	//	Input 249 Output 999
	//	Input 250 Output 1003
	//	Input 251 Output 1007
	//	Input 252 Output 1011
	//	Input 253 Output 1015
	//	Input 254 Output 1019
	//	Input 255 Output 1023

}

func ExampleRangeNegative255() {

	a1 := int64(0)
	a2 := int64(-255)
	b1 := int64(0)
	b2 := int64(1023)


	for i:=int64(0);i>-256;i-- {
		r := NewRangeMapFrom(a1, a2, b1, b2)
		fmt.Printf("\tInput %d Output %d\n",i, r.MustMap(i))
	}

	// Output:
//	Input 0 Output 0
//	Input -1 Output 4
//	Input -2 Output 8
//	Input -3 Output 12
//	Input -4 Output 16
//	Input -5 Output 20
//	Input -6 Output 24
//	Input -7 Output 28
//	Input -8 Output 32
//	Input -9 Output 36
//	Input -10 Output 40
//	Input -11 Output 44
//	Input -12 Output 48
//	Input -13 Output 52
//	Input -14 Output 56
//	Input -15 Output 60
//	Input -16 Output 64
//	Input -17 Output 68
//	Input -18 Output 72
//	Input -19 Output 76
//	Input -20 Output 80
//	Input -21 Output 84
//	Input -22 Output 88
//	Input -23 Output 92
//	Input -24 Output 96
//	Input -25 Output 100
//	Input -26 Output 104
//	Input -27 Output 108
//	Input -28 Output 112
//	Input -29 Output 116
//	Input -30 Output 120
//	Input -31 Output 124
//	Input -32 Output 128
//	Input -33 Output 132
//	Input -34 Output 136
//	Input -35 Output 140
//	Input -36 Output 144
//	Input -37 Output 148
//	Input -38 Output 152
//	Input -39 Output 156
//	Input -40 Output 160
//	Input -41 Output 164
//	Input -42 Output 168
//	Input -43 Output 173
//	Input -44 Output 177
//	Input -45 Output 181
//	Input -46 Output 185
//	Input -47 Output 189
//	Input -48 Output 193
//	Input -49 Output 197
//	Input -50 Output 201
//	Input -51 Output 205
//	Input -52 Output 209
//	Input -53 Output 213
//	Input -54 Output 217
//	Input -55 Output 221
//	Input -56 Output 225
//	Input -57 Output 229
//	Input -58 Output 233
//	Input -59 Output 237
//	Input -60 Output 241
//	Input -61 Output 245
//	Input -62 Output 249
//	Input -63 Output 253
//	Input -64 Output 257
//	Input -65 Output 261
//	Input -66 Output 265
//	Input -67 Output 269
//	Input -68 Output 273
//	Input -69 Output 277
//	Input -70 Output 281
//	Input -71 Output 285
//	Input -72 Output 289
//	Input -73 Output 293
//	Input -74 Output 297
//	Input -75 Output 301
//	Input -76 Output 305
//	Input -77 Output 309
//	Input -78 Output 313
//	Input -79 Output 317
//	Input -80 Output 321
//	Input -81 Output 325
//	Input -82 Output 329
//	Input -83 Output 333
//	Input -84 Output 337
//	Input -85 Output 341
//	Input -86 Output 345
//	Input -87 Output 349
//	Input -88 Output 353
//	Input -89 Output 357
//	Input -90 Output 361
//	Input -91 Output 365
//	Input -92 Output 369
//	Input -93 Output 373
//	Input -94 Output 377
//	Input -95 Output 381
//	Input -96 Output 385
//	Input -97 Output 389
//	Input -98 Output 393
//	Input -99 Output 397
//	Input -100 Output 401
//	Input -101 Output 405
//	Input -102 Output 409
//	Input -103 Output 413
//	Input -104 Output 417
//	Input -105 Output 421
//	Input -106 Output 425
//	Input -107 Output 429
//	Input -108 Output 433
//	Input -109 Output 437
//	Input -110 Output 441
//	Input -111 Output 445
//	Input -112 Output 449
//	Input -113 Output 453
//	Input -114 Output 457
//	Input -115 Output 461
//	Input -116 Output 465
//	Input -117 Output 469
//	Input -118 Output 473
//	Input -119 Output 477
//	Input -120 Output 481
//	Input -121 Output 485
//	Input -122 Output 489
//	Input -123 Output 493
//	Input -124 Output 497
//	Input -125 Output 501
//	Input -126 Output 505
//	Input -127 Output 509
//	Input -128 Output 514
//	Input -129 Output 518
//	Input -130 Output 522
//	Input -131 Output 526
//	Input -132 Output 530
//	Input -133 Output 534
//	Input -134 Output 538
//	Input -135 Output 542
//	Input -136 Output 546
//	Input -137 Output 550
//	Input -138 Output 554
//	Input -139 Output 558
//	Input -140 Output 562
//	Input -141 Output 566
//	Input -142 Output 570
//	Input -143 Output 574
//	Input -144 Output 578
//	Input -145 Output 582
//	Input -146 Output 586
//	Input -147 Output 590
//	Input -148 Output 594
//	Input -149 Output 598
//	Input -150 Output 602
//	Input -151 Output 606
//	Input -152 Output 610
//	Input -153 Output 614
//	Input -154 Output 618
//	Input -155 Output 622
//	Input -156 Output 626
//	Input -157 Output 630
//	Input -158 Output 634
//	Input -159 Output 638
//	Input -160 Output 642
//	Input -161 Output 646
//	Input -162 Output 650
//	Input -163 Output 654
//	Input -164 Output 658
//	Input -165 Output 662
//	Input -166 Output 666
//	Input -167 Output 670
//	Input -168 Output 674
//	Input -169 Output 678
//	Input -170 Output 682
//	Input -171 Output 686
//	Input -172 Output 690
//	Input -173 Output 694
//	Input -174 Output 698
//	Input -175 Output 702
//	Input -176 Output 706
//	Input -177 Output 710
//	Input -178 Output 714
//	Input -179 Output 718
//	Input -180 Output 722
//	Input -181 Output 726
//	Input -182 Output 730
//	Input -183 Output 734
//	Input -184 Output 738
//	Input -185 Output 742
//	Input -186 Output 746
//	Input -187 Output 750
//	Input -188 Output 754
//	Input -189 Output 758
//	Input -190 Output 762
//	Input -191 Output 766
//	Input -192 Output 770
//	Input -193 Output 774
//	Input -194 Output 778
//	Input -195 Output 782
//	Input -196 Output 786
//	Input -197 Output 790
//	Input -198 Output 794
//	Input -199 Output 798
//	Input -200 Output 802
//	Input -201 Output 806
//	Input -202 Output 810
//	Input -203 Output 814
//	Input -204 Output 818
//	Input -205 Output 822
//	Input -206 Output 826
//	Input -207 Output 830
//	Input -208 Output 834
//	Input -209 Output 838
//	Input -210 Output 842
//	Input -211 Output 846
//	Input -212 Output 850
//	Input -213 Output 855
//	Input -214 Output 859
//	Input -215 Output 863
//	Input -216 Output 867
//	Input -217 Output 871
//	Input -218 Output 875
//	Input -219 Output 879
//	Input -220 Output 883
//	Input -221 Output 887
//	Input -222 Output 891
//	Input -223 Output 895
//	Input -224 Output 899
//	Input -225 Output 903
//	Input -226 Output 907
//	Input -227 Output 911
//	Input -228 Output 915
//	Input -229 Output 919
//	Input -230 Output 923
//	Input -231 Output 927
//	Input -232 Output 931
//	Input -233 Output 935
//	Input -234 Output 939
//	Input -235 Output 943
//	Input -236 Output 947
//	Input -237 Output 951
//	Input -238 Output 955
//	Input -239 Output 959
//	Input -240 Output 963
//	Input -241 Output 967
//	Input -242 Output 971
//	Input -243 Output 975
//	Input -244 Output 979
//	Input -245 Output 983
//	Input -246 Output 987
//	Input -247 Output 991
//	Input -248 Output 995
//	Input -249 Output 999
//	Input -250 Output 1003
//	Input -251 Output 1007
//	Input -252 Output 1011
//	Input -253 Output 1015
//	Input -254 Output 1019
//	Input -255 Output 1023

}


func ExampleRangeBoolean255() {

	a1 := int64(0)
	a2 := int64(255)
	b1 := int64(0)
	b2 := int64(1)

	for i:=int64(0);i<256;i++ {
		r := NewRangeMapFrom(a1, a2, b1, b2)
		fmt.Printf("\tInput %d Output %d\n",i, r.MustMap(i))
	}
// Output:
//	Input 0 Output 0
//	Input 1 Output 0
//	Input 2 Output 0
//	Input 3 Output 0
//	Input 4 Output 0
//	Input 5 Output 0
//	Input 6 Output 0
//	Input 7 Output 0
//	Input 8 Output 0
//	Input 9 Output 0
//	Input 10 Output 0
//	Input 11 Output 0
//	Input 12 Output 0
//	Input 13 Output 0
//	Input 14 Output 0
//	Input 15 Output 0
//	Input 16 Output 0
//	Input 17 Output 0
//	Input 18 Output 0
//	Input 19 Output 0
//	Input 20 Output 0
//	Input 21 Output 0
//	Input 22 Output 0
//	Input 23 Output 0
//	Input 24 Output 0
//	Input 25 Output 0
//	Input 26 Output 0
//	Input 27 Output 0
//	Input 28 Output 0
//	Input 29 Output 0
//	Input 30 Output 0
//	Input 31 Output 0
//	Input 32 Output 0
//	Input 33 Output 0
//	Input 34 Output 0
//	Input 35 Output 0
//	Input 36 Output 0
//	Input 37 Output 0
//	Input 38 Output 0
//	Input 39 Output 0
//	Input 40 Output 0
//	Input 41 Output 0
//	Input 42 Output 0
//	Input 43 Output 0
//	Input 44 Output 0
//	Input 45 Output 0
//	Input 46 Output 0
//	Input 47 Output 0
//	Input 48 Output 0
//	Input 49 Output 0
//	Input 50 Output 0
//	Input 51 Output 0
//	Input 52 Output 0
//	Input 53 Output 0
//	Input 54 Output 0
//	Input 55 Output 0
//	Input 56 Output 0
//	Input 57 Output 0
//	Input 58 Output 0
//	Input 59 Output 0
//	Input 60 Output 0
//	Input 61 Output 0
//	Input 62 Output 0
//	Input 63 Output 0
//	Input 64 Output 0
//	Input 65 Output 0
//	Input 66 Output 0
//	Input 67 Output 0
//	Input 68 Output 0
//	Input 69 Output 0
//	Input 70 Output 0
//	Input 71 Output 0
//	Input 72 Output 0
//	Input 73 Output 0
//	Input 74 Output 0
//	Input 75 Output 0
//	Input 76 Output 0
//	Input 77 Output 0
//	Input 78 Output 0
//	Input 79 Output 0
//	Input 80 Output 0
//	Input 81 Output 0
//	Input 82 Output 0
//	Input 83 Output 0
//	Input 84 Output 0
//	Input 85 Output 0
//	Input 86 Output 0
//	Input 87 Output 0
//	Input 88 Output 0
//	Input 89 Output 0
//	Input 90 Output 0
//	Input 91 Output 0
//	Input 92 Output 0
//	Input 93 Output 0
//	Input 94 Output 0
//	Input 95 Output 0
//	Input 96 Output 0
//	Input 97 Output 0
//	Input 98 Output 0
//	Input 99 Output 0
//	Input 100 Output 0
//	Input 101 Output 0
//	Input 102 Output 0
//	Input 103 Output 0
//	Input 104 Output 0
//	Input 105 Output 0
//	Input 106 Output 0
//	Input 107 Output 0
//	Input 108 Output 0
//	Input 109 Output 0
//	Input 110 Output 0
//	Input 111 Output 0
//	Input 112 Output 0
//	Input 113 Output 0
//	Input 114 Output 0
//	Input 115 Output 0
//	Input 116 Output 0
//	Input 117 Output 0
//	Input 118 Output 0
//	Input 119 Output 0
//	Input 120 Output 0
//	Input 121 Output 0
//	Input 122 Output 0
//	Input 123 Output 0
//	Input 124 Output 0
//	Input 125 Output 0
//	Input 126 Output 0
//	Input 127 Output 0
//	Input 128 Output 1
//	Input 129 Output 1
//	Input 130 Output 1
//	Input 131 Output 1
//	Input 132 Output 1
//	Input 133 Output 1
//	Input 134 Output 1
//	Input 135 Output 1
//	Input 136 Output 1
//	Input 137 Output 1
//	Input 138 Output 1
//	Input 139 Output 1
//	Input 140 Output 1
//	Input 141 Output 1
//	Input 142 Output 1
//	Input 143 Output 1
//	Input 144 Output 1
//	Input 145 Output 1
//	Input 146 Output 1
//	Input 147 Output 1
//	Input 148 Output 1
//	Input 149 Output 1
//	Input 150 Output 1
//	Input 151 Output 1
//	Input 152 Output 1
//	Input 153 Output 1
//	Input 154 Output 1
//	Input 155 Output 1
//	Input 156 Output 1
//	Input 157 Output 1
//	Input 158 Output 1
//	Input 159 Output 1
//	Input 160 Output 1
//	Input 161 Output 1
//	Input 162 Output 1
//	Input 163 Output 1
//	Input 164 Output 1
//	Input 165 Output 1
//	Input 166 Output 1
//	Input 167 Output 1
//	Input 168 Output 1
//	Input 169 Output 1
//	Input 170 Output 1
//	Input 171 Output 1
//	Input 172 Output 1
//	Input 173 Output 1
//	Input 174 Output 1
//	Input 175 Output 1
//	Input 176 Output 1
//	Input 177 Output 1
//	Input 178 Output 1
//	Input 179 Output 1
//	Input 180 Output 1
//	Input 181 Output 1
//	Input 182 Output 1
//	Input 183 Output 1
//	Input 184 Output 1
//	Input 185 Output 1
//	Input 186 Output 1
//	Input 187 Output 1
//	Input 188 Output 1
//	Input 189 Output 1
//	Input 190 Output 1
//	Input 191 Output 1
//	Input 192 Output 1
//	Input 193 Output 1
//	Input 194 Output 1
//	Input 195 Output 1
//	Input 196 Output 1
//	Input 197 Output 1
//	Input 198 Output 1
//	Input 199 Output 1
//	Input 200 Output 1
//	Input 201 Output 1
//	Input 202 Output 1
//	Input 203 Output 1
//	Input 204 Output 1
//	Input 205 Output 1
//	Input 206 Output 1
//	Input 207 Output 1
//	Input 208 Output 1
//	Input 209 Output 1
//	Input 210 Output 1
//	Input 211 Output 1
//	Input 212 Output 1
//	Input 213 Output 1
//	Input 214 Output 1
//	Input 215 Output 1
//	Input 216 Output 1
//	Input 217 Output 1
//	Input 218 Output 1
//	Input 219 Output 1
//	Input 220 Output 1
//	Input 221 Output 1
//	Input 222 Output 1
//	Input 223 Output 1
//	Input 224 Output 1
//	Input 225 Output 1
//	Input 226 Output 1
//	Input 227 Output 1
//	Input 228 Output 1
//	Input 229 Output 1
//	Input 230 Output 1
//	Input 231 Output 1
//	Input 232 Output 1
//	Input 233 Output 1
//	Input 234 Output 1
//	Input 235 Output 1
//	Input 236 Output 1
//	Input 237 Output 1
//	Input 238 Output 1
//	Input 239 Output 1
//	Input 240 Output 1
//	Input 241 Output 1
//	Input 242 Output 1
//	Input 243 Output 1
//	Input 244 Output 1
//	Input 245 Output 1
//	Input 246 Output 1
//	Input 247 Output 1
//	Input 248 Output 1
//	Input 249 Output 1
//	Input 250 Output 1
//	Input 251 Output 1
//	Input 252 Output 1
//	Input 253 Output 1
//	Input 254 Output 1
//	Input 255 Output 1




}



func ExampleRangeOver() {

	a1 := int64(0)
	a2 := int64(255)
	b1 := int64(0)
	b2 := int64(250)

	i := int64(500)
	r := NewRangeMapFrom(a1, a2, b1, b2)
	fmt.Printf("\tInput %d Output %d\n",i, r.MustMap(i))

// Output:
//	Input 500 Output 250

}

func ExampleRangeUnder() {

	a1 := int64(0)
	a2 := int64(255)
	b1 := int64(0)
	b2 := int64(250)

	i := int64(-500)
	r := NewRangeMapFrom(a1, a2, b1, b2)
	fmt.Printf("\tInput %d Output %d\n",i, r.MustMap(i))

// Output:
//	Input -500 Output 0

}


func ExampleRangeToBool() {

	a1 := int64(255)
	a2 := int64(170)
	b1 := int64(0)
	b2 := int64(1)

	i := int64(212)
	r := NewRangeMapFrom(a1, a2, b1, b2)
	fmt.Printf("\tInput %d Output %d\n",i, r.MustMap(i))
	fmt.Printf("\tInput %d Output %d\n",i+2, r.MustMap(i+2))

	// Output:
	//	Input 212 Output 1
	//	Input 214 Output 0

}
