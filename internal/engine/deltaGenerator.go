package engine

// this file will generate the delta for converting version x to version x+ 1

// this files received the LCS array , the Xth version array and the X+1th version array
// generate the smallest set of opeation needed to be done on file X to convert it to X + 1 th file

// return this set of operation

func GenerateDelta(versionX1, versionX2, lcs *[]int) []DeltaInstruction {

	var delta []DeltaInstruction
	idx1 := 0
	idx2 := 0
	idxLCS := 0

	for idx1 < len((*versionX1)) || idx2 < len((*versionX2)) || idxLCS < len((*lcs)) {

		// if LCS ended and parrent remaining

		if idx1 < len((*versionX1)) && idxLCS == len((*lcs)) {
			delta = append(delta, DeltaInstruction{
				DeltaType: "D",
				Line:      idx1,
				Val:       0, // 0 just measa garbage value never use it to identify a string
			})
			idx1++
		} else if idx2 < len((*versionX2)) && idxLCS == len((*lcs)) {
			// verxionX2 remaining
			delta = append(delta, DeltaInstruction{
				DeltaType: "A",
				Line:      idx1,
				Val:       (*versionX2)[idx2],
			})
			idx2++
		} else {

			// idxLCS != len(lcs)

			// case 1 al 3 matches

			if (*versionX1)[idx1] == (*lcs)[idxLCS] && (*versionX2)[idx2] == (*lcs)[idxLCS] {
				idx1++
				idx2++
				idxLCS++
			} else {

				// case 2 versionX does not match with LCS Element

				if (*versionX2)[idx2] != (*lcs)[idxLCS] {
					delta = append(delta, DeltaInstruction{
						DeltaType: "A",
						Line:      idx1,
						Val:       (*versionX2)[idx2],
					})
					idx2++
				} else if (*versionX1)[idx1] != (*lcs)[idxLCS] {
					delta = append(delta, DeltaInstruction{
						DeltaType: "D",
						Line:      idx1,
						Val:       0,
					})
					idx1++
				}

				// case 3 versionX does not match with LCS Element
			}

		}
	}

	return delta

}
