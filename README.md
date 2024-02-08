***Instructions***

1) Execute the command "go run No_Concurrency.go".  This will run the go program version without concurrency and print out the average MSE, AIC, and BIC values from 100 sequential runs to the terminal.

2) Execute the command "go run Concurrency.go".  This will run the go program version with concurrency and print out the average MSE, AIC, and BIC values from 100 sequential runs to the terminal.

 ***Notes to Management***

 I was unable to successfully find a way to run every possible linear regression. Accordingly, I only ran a simple linear regression - albeit 100 times.  The results showed that the version without concurrency actually ran faster.  My assumption is that if I was building every possible linear regression that I would have much larger processing time and the concurrency version would execute much more quickly.   