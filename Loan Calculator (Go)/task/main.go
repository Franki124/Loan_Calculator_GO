package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
)

func main() {
	typeFlag := flag.String("type", "", "Type of the payment calculation: 'annuity' or 'diff'")
	principalFlag := flag.String("principal", "", "The loan principal")
	periodsFlag := flag.String("periods", "", "The number of periods (months) to repay the loan")
	interestFlag := flag.String("interest", "", "The loan interest rate as a percentage without the % sign")
	paymentFlag := flag.String("payment", "", "The payment amount per period")

	flag.Parse()

	if *interestFlag == "" {
		fmt.Println("Incorrect parameters: Interest rate must be provided.")
		return
	}

	principal, errPrincipal := strconv.ParseFloat(*principalFlag, 64)
	periods, errPeriods := strconv.Atoi(*periodsFlag)
	interest, errInterest := strconv.ParseFloat(*interestFlag, 64)
	payment, errPayment := strconv.ParseFloat(*paymentFlag, 64)

	if errInterest != nil || interest <= 0 || (errPrincipal == nil && principal <= 0) || (errPeriods == nil && periods <= 0) || (errPayment == nil && payment <= 0) {
		fmt.Println("Incorrect parameters")
		return
	}

	monthlyInterest := interest / 12 / 100

	if *typeFlag == "annuity" {
		if *paymentFlag == "" && errPrincipal == nil && errPeriods == nil {
			payment = principal * (monthlyInterest * math.Pow(1+monthlyInterest, float64(periods))) / (math.Pow(1+monthlyInterest, float64(periods)) - 1)
			payment = math.Ceil(payment)
			totalPayment := payment * float64(periods)
			overpayment := totalPayment - principal
			fmt.Printf("Your monthly payment = $%d!\n", int(payment))
			fmt.Printf("Overpayment = $%d\n", int(overpayment))
		} else if *principalFlag == "" && errPayment == nil && errPeriods == nil {
			principal = payment / ((monthlyInterest * math.Pow(1+monthlyInterest, float64(periods))) / (math.Pow(1+monthlyInterest, float64(periods)) - 1))
			overpayment := float64(periods)*payment - principal
			fmt.Printf("Your loan principal = $%d\n", int(principal))
			fmt.Printf("Overpayment = $%d\n", int(overpayment))
		} else if *periodsFlag == "" && errPrincipal == nil && errPayment == nil {
			n := math.Log(payment/(payment-monthlyInterest*principal)) / math.Log(1+monthlyInterest)
			months := int(math.Ceil(n))
			totalPayment := payment * float64(months)
			overpayment := totalPayment - principal

			years := months / 12
			months %= 12
			if years > 0 {
				if months > 0 {
					fmt.Printf("It will take %d years and %d months to repay this loan!\n", years, months)
				} else {
					fmt.Printf("It will take %d years to repay this loan!\n", years)
				}
			} else {
				fmt.Printf("It will take %d months to repay this loan!\n", months)
			}
			fmt.Printf("Overpayment = $%d\n", int(overpayment))
		} else {
			fmt.Println("Incorrect parameters")
		}
	} else if *typeFlag == "diff" {
		if *paymentFlag != "" {
			fmt.Println("Incorrect parameters: Payment parameter should not be provided for 'diff' type.")
			return
		}
		if errPrincipal != nil || principal <= 0 || errPeriods != nil || periods <= 0 {
			fmt.Println("Incorrect parameters")
			return
		}
		totalPayment := 0.0
		for m := 1; m <= periods; m++ {
			dm := principal/float64(periods) + monthlyInterest*(principal-principal*float64(m-1)/float64(periods))
			dm = math.Ceil(dm)
			fmt.Printf("Month %d: payment is $%d\n", m, int(dm))
			totalPayment += dm
		}
		overpayment := totalPayment - principal
		fmt.Printf("Overpayment = $%d\n", int(math.Ceil(overpayment)))
	} else {
		fmt.Println("Incorrect parameters: Type must be either 'annuity' or 'diff'.")
	}
}
