Explain what the following code is attempting to do? You can explain by:
1. Explaining how the highlighted constructs work?
2. Giving use-cases of what these constructs could be used for.
3. What is the significance of the for loop with 4 iterations?
4. What is the significance of make(chan func(), 10)?
5. Why is “HERE1” not getting printed?


```package main
import &quot;fmt&quot;
func main() {
cnp := make(chan func(), 10)
for i := 0; i &lt; 4; i++ {
go func() {
for f := range cnp {
f()
}
}()
}
cnp &lt;- func() {
fmt.Println(&quot;HERE1&quot;)
}
fmt.Println(&quot;Hello&quot;)
}```