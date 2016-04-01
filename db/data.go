package db

type Protein struct {
	ID string
	//original
	Seq    []string
	Dssp   []string
	Stride []string
	Kaksi  []string
	Pross  []string
	//processed
	Dssp3   []string
	Stride3 []string
	Kaksi3  []string
	Pross3  []string
	// consensus 2
	DsspStride3  []string
	DsspKaksi3   []string
	DsspPross3   []string
	StrideKaksi3 []string
	StridePross3 []string
	KaksiPross3  []string
	// consensus 3
	DsspStrideKaksi3  []string
	DsspStridePross3  []string
	DsspKaksiPross3   []string
	StrideKaksiPross3 []string
	// consensus 4
	All3 []string
	// hydrophobicity rose
	SeqHPRose               []string
	DsspHPRose              []string
	StrideHPRose            []string
	KaksiHPRose             []string
	ProssHPRose             []string
	Dssp3HPRose             []string
	Stride3HPRose           []string
	Kaksi3HPRose            []string
	Pross3HPRose            []string
	DsspStride3HPRose       []string
	DsspKaksi3HPRose        []string
	DsspPross3HPRose        []string
	StrideKaksi3HPRose      []string
	StridePross3HPRose      []string
	KaksiPross3HPRose       []string
	DsspStrideKaksi3HPRose  []string
	DsspStridePross3HPRose  []string
	DsspKaksiPross3HPRose   []string
	StrideKaksiPross3HPRose []string
	All3HPRose              []string
	// rose special: polar, nonpolar, Gly and Pro
	SeqHPRoseSpecial               []string
	DsspHPRoseSpecial              []string
	StrideHPRoseSpecial            []string
	KaksiHPRoseSpecial             []string
	ProssHPRoseSpecial             []string
	Dssp3HPRoseSpecial             []string
	Stride3HPRoseSpecial           []string
	Kaksi3HPRoseSpecial            []string
	Pross3HPRoseSpecial            []string
	DsspStride3HPRoseSpecial       []string
	DsspKaksi3HPRoseSpecial        []string
	DsspPross3HPRoseSpecial        []string
	StrideKaksi3HPRoseSpecial      []string
	StridePross3HPRoseSpecial      []string
	KaksiPross3HPRoseSpecial       []string
	DsspStrideKaksi3HPRoseSpecial  []string
	DsspStridePross3HPRoseSpecial  []string
	DsspKaksiPross3HPRoseSpecial   []string
	StrideKaksiPross3HPRoseSpecial []string
	All3HPRoseSpecial              []string
	//rose special charged: polar, nonpolar, Gly, Pro, positives and negatives
	SeqHPRoseSpecialCharged               []string
	DsspHPRoseSpecialCharged              []string
	StrideHPRoseSpecialCharged            []string
	KaksiHPRoseSpecialCharged             []string
	ProssHPRoseSpecialCharged             []string
	Dssp3HPRoseSpecialCharged             []string
	Stride3HPRoseSpecialCharged           []string
	Kaksi3HPRoseSpecialCharged            []string
	Pross3HPRoseSpecialCharged            []string
	DsspStride3HPRoseSpecialCharged       []string
	DsspKaksi3HPRoseSpecialCharged        []string
	DsspPross3HPRoseSpecialCharged        []string
	StrideKaksi3HPRoseSpecialCharged      []string
	StridePross3HPRoseSpecialCharged      []string
	KaksiPross3HPRoseSpecialCharged       []string
	DsspStrideKaksi3HPRoseSpecialCharged  []string
	DsspStridePross3HPRoseSpecialCharged  []string
	DsspKaksiPross3HPRoseSpecialCharged   []string
	StrideKaksiPross3HPRoseSpecialCharged []string
	All3HPRoseSpecialCharged              []string
}
