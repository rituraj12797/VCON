package engine

import (
	"fmt"
	"time"
	"vcon/internal/globalStore"
)

// this will convert the identifier list into readable content

func ContentRendered(contentNumArray *[]int) {

	size := len(*contentNumArray)

	resultArray := make([]string, size)

	fmt.Println(" Received array : ", contentNumArray)
	// use GetStringFromIdentifier for getting string from identifier
	gStore := globalStore.GlobalStore

	startTime := time.Now() // Record start time

	for i, y := range *contentNumArray {
		// y is the identifier
		var stringForThisIdentifier, err = gStore.GetStringFromIdentifier(y)
		if err != nil {
			// Error happeed here
			fmt.Println("Error happened : ", err)
		} else {
			// string parsed succesfully
			resultArray[i] = stringForThisIdentifier // stored sequentially order maintained dueto iteration
		}
	}
	
	/*



			NICE ENGINEERING BUT OVERKILL KILLED PERFORMANCE : the computation time is very less then the time overhead due to multiple go rouines and channels so theretically the concurent approach should have been faster it turned out to be slow even uptill 10M lines of text file



			if size < 8 {
				// fetch single threadedly ( using the main thread ) and fill the resultArray

				for i, y := range *contentNumArray {
					// y is the identifier
					var stringForThisIdentifier, err = gStore.GetStringFromIdentifier(y)

					if err != nil {
						// Error happeed here
						fmt.Println("Error happened : ", err)
					} else {
						// string parsed succesfully
						resultArray[i] = stringForThisIdentifier // stored sequentially order maintained due to iteration
					}
				}
			} else {
				// fetch using multiple (8) threads and each will write the result to specifiic corrospondign indec in the result array to main order

				// add 8 go routines
				// initially each will read from gStore
				// basically each go rouine will be responsible for a range of indices not a single indice

				// ranges will  be like [0,n/8 + 1],[n/8+2,2n/8+3],....

				// make a vector of pairs // vector [i] = start_i, end_i

				// last time this ued a go routine to manage 8 lined chunks which is eneffcient and used upto 125 go routines in 1000 lined texts need to correct it to correclty leverage concurrecy


				var rangeArr []PairRange
		        numWorkers := 16
		        // this calculation correctly determines how many items each of the 8 workers should process.
		        // it uses ceiling division to handle remainders gracefully.
		        chunkSize := (size + numWorkers - 1) / numWorkers

		        for i := 0; i < numWorkers; i++ {
		            start := i * chunkSize
		            // the end index is inclusive for your loop `x <= endIndex`
		            end := start + chunkSize - 1

		            // ff a worker's calculated start is beyond the end of the slice, it has no work.
		            if start >= size {
		                continue
		            }

		            // The last worker's end index must be clamped to the actual end of the slice.
		            if end >= size {
		                end = size - 1
		            }

		            rangeArr = append(rangeArr, PairRange{
		                start: start,
		                end:   end,
		            })
		        }

				// ranges defined
				// fmt.Println(" Ranges :  ", rangeArr)

				var wg sync.WaitGroup

				type resultPayload struct {
					index     int
					statement string
				}

				ch := make(chan resultPayload, size) // only this number of contnts will come

				for _, curRange := range rangeArr {
					wg.Add(1)

					go func(startIndex int, endIndex int) {
						defer wg.Done()

						for x := startIndex; x <= endIndex; x++ {
							var stringForThisIdentifier, err = gStore.GetStringFromIdentifier((*contentNumArray)[x])
							var statementToPush string

							// fmt.Println(" ID : ", (*contentNumArray)[x], "   string: ", stringForThisIdentifier)

							if err != nil {
								statementToPush = fmt.Sprintf("error id %d not found (%v)", (*contentNumArray)[x], err)
							} else {
								statementToPush = stringForThisIdentifier
							}

							ch <- resultPayload{
								index:     x,
								statement: statementToPush,
							}
						}

					}(curRange.start, curRange.end)
				}

				go func() {

					// to close the channel ==> WHY ??

						*LEARNING : The loop for consumer reading from the chnnel runs concurrently with producer pushing in the content
						if the channle is empty the range loop will be blocked and wait untill something comes into the channle
						now sa if all go routines are done reading from the store and now nothing will come into the channle
						the consumer will be blocked forever and DEADLOCK will happen to prevent this we do add a era go routine whihc will wait untill the wait group is done

						WG done => all go routines executed => nothign will come into channel so close


					wg.Wait()
					close(ch)
				}()

				// main thread will read from channel

				for payload := range ch { // payloas is a incoming object from the channle and the loop will beblocked untill the channle is empty
					resultArray[payload.index] = payload.statement
				}
			}
	*/

	duration := time.Since(startTime) // Calculate duration
	
	// 5. Print the result
	fmt.Printf("\n--- Benchmark Complete ---\n")

	// print the result now
	fmt.Println(" =========== COPMILED RESULT ========== ")

	for x, y := range resultArray {
		fmt.Println(x, " :  ", y)
	}

	fmt.Println(" ====================================== ")
	fmt.Printf("ContentRenderer took %v micro seconds to process %d lines.\n", duration.Microseconds(), size)

}
