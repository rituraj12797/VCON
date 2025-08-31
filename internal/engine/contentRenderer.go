package engine

import (
	"fmt"
	"sync"
	"vcon/internal/globalStore"
)

// this will convert the identifier list into readable content


type PairRange struct {
	start int;
	end int;
}

func ContentRendered(contentNumArray *[]int) {

	size := len(*contentNumArray);

	resultArray := make([]string, size);
	
	
	// use GetStringFromIdentifier for getting string from identifier 
	gStore := globalStore.GlobalStore

	if(size < 8) {
		// fetch single threadedly ( using the main thread ) and fill the resultArray

		for i,y := range *contentNumArray {
			// y is the identifier 
			var stringForThisIdentifier, err = gStore.GetStringFromIdentifier(y)
			
			if err != nil {
				// Error happeed here
				fmt.Println("Error happened : ",err);
			} else {
				// string parsed succesfully 
				resultArray[i] = stringForThisIdentifier; // stored sequentially order maintained due to iteration
			}
		}
	} else {
		// fetch using multiple (8) threads and each will write the result to specifiic corrospondign indec in the result array to main order 


		// add 8 go routines 
		// initially each will read from gStore 
		// basically each go rouine will be responsible for a range of indices not a single indice

		// ranges will  be like [0,n/8 + 1],[n/8+2,2n/8+3],....

		// make a vector of pairs // vector [i] = start_i, end_i

		var rangeArr []PairRange;

		for i:= 0; i < size; i+=8 {
			rangeArr = append(rangeArr, PairRange{
				start: i,
				end: min(i+7,size - 1),
			})
		}

		// ranges defined 
		fmt.Println(" Ranges :  ", rangeArr)


		var wg sync.WaitGroup;

		type resultPayload struct  {
			index int;
			statement string;
		}

		ch := make(chan resultPayload, size ) // only this number of contnts will come 
		

		for _,curRange := range rangeArr{
			wg.Add(1)

			go func( startIndex int, endIndex int) {
				defer wg.Done()

				for x := startIndex ; x <= endIndex; x++ {
					var stringForThisIdentifier, err = gStore.GetStringFromIdentifier((*contentNumArray)[x])
					var statementToPush string;
					if err != nil {
						statementToPush = fmt.Sprintf("error id %d not found (%v)", (*contentNumArray)[x], err)
					} else {
						statementToPush = stringForThisIdentifier
					}

					ch<-resultPayload{
						index : x,
						statement: statementToPush,
					}
				}

			}(curRange.start,curRange.end)
		}

		


		go func() {

			// to close the channel ==> WHY ?? 
		/*

		*LEARNING : The loop for consumer reading from the chnnel runs concurrently with producer pushing in the content
		if the channle is empty the range loop will be blocked and wait untill something comes into the channle 
		now sa if all go routines are done reading from the store and now nothing will come into the channle 
		the consumer will be blocked forever and DEADLOCK will happen to prevent this we do add a era go routine whihc will wait untill the wait group is done
		
		WG done => all go routines executed => nothign will come into channel so close 
		*/

			wg.Wait()
			close(ch)
		}()


		// main thread will read from channel 

		for payload := range ch { // payloas is a incoming object from the channle and the loop will beblocked untill the channle is empty 
			resultArray[payload.index] = payload.statement
		}
	}

	// print the result now 
	fmt.Println(" =========== COPMILED RESULT ========== ")
	
	for x,y := range resultArray {
		fmt.Println(x," :  ",y);
	}

	fmt.Println(" ====================================== ")
	
	

}