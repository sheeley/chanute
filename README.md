# Chanute
A library that makes dealing with AWS Trusted Advisor less bad.

## Who is Chanute?
[Octave Chanute](https://en.wikipedia.org/wiki/Octave_Chanute) was the Wright Brothers' trusted advisor and the father of aviation.

## Examples
View the `cmd` directory for more complex usage.

```
sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
r, err := chanute.GenerateReport(sess)
if err != nil {
    panic(err)
}
fmt.Println(r.AsciiReport())
```

### Outputs
```
EC2
+-----------------------------+----+----------------------+---------------------------+
|            NAME             | ID | LOW UTILIZATION DAYS | ESTIMATED MONTHLY SAVINGS |
+-----------------------------+----+----------------------+---------------------------+
| cloud engineering           |    |                      | $2,792                    |
| sct                         |    |                      | $303                      |
+-----------------------------+----+----------------------+---------------------------+

Load Balancers
+----------------------------------+--------+--------+--------------+
|               NAME               | REGION | REASON | MONTHLY COST |
+----------------------------------+--------+--------+--------------+
| Production-Backend-LoadBalancer  |        |        | $18          |
+----------------------------------+--------+--------+--------------+

EBS
+-------------------------------------------------------------+-------------+--------------+--------------+
|                          VOLUME ID                          | VOLUME NAME | SIZE (IN GB) | MONTHLY COST |
+-------------------------------------------------------------+-------------+--------------+--------------+
| site reliability                                            |             |          112 | $10          |
+-------------------------------------------------------------+-------------+--------------+--------------+

RDS
+----------------------+---------+-----------------------+----------------------+--------------+
|         NAME         | MULTIAZ | DAYS SINCE CONNECTION | STORAGE SIZE (IN GB) | MONTHLY COST |
+----------------------+---------+-----------------------+----------------------+--------------+
| cloud engineering    |         |                       |                    2 | $60          |
+----------------------+---------+-----------------------+----------------------+--------------+

Redshift
+-------------------+--------+--------+--------------+
|       NAME        | STATUS | REASON | MONTHLY COST |
+-------------------+--------+--------+--------------+
| cloud engineering |        |        | $180         |
+-------------------+--------+--------+--------------+
```