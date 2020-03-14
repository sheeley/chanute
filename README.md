# Chanute
A library that makes dealing with AWS Trusted Advisor less bad.

## Who is Chanute?
[Octave Chanute](https://en.wikipedia.org/wiki/Octave_Chanute) was the Wright Brothers' trusted advisor and the father of aviation.

## Examples
```
sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
r, err := chanute.GenerateReport(sess)
if err != nil {
    panic(err)
}
fmt.Println(r.AsciiReport())
```

View the `cmd` directory for more.