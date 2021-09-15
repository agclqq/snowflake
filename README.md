# snowflake
[中文说明](README_ZH.md)

This is the snowflake algorithm of a GO language implementation.

Support for more flexible generation methods.

## goal：
1. In a distributed environment, the sequence number is unique
2. In a multithreading environment,the sequence number is unique
3. Standalone performance is good enough

## implement solution：
The generated serial number is a 64-bit number.In order to achieve this goal, this number is made up of the following components:
* datacenter
* machine
* timestamp
* sequence

Datacenter and machine to achieve multi-machine cases, each machine has a unique identifier;

Timestamps ensure orderly growth and avoid repetition in the case of a single machine;

The sequence is guaranteed to increase in unit time;

In addition, the program must ensure that the sequence generation of multiple processes is not repeated;
## Sequence number resolution：
The binary structure of the serial number：
```go
0 0000 0000 0000000000000000000000000000000000000000 000000000000000
```
Group by space：
* The first group，The value is fixed to 0 to avoid data overflow and negative values
* The second group，The length is not fixed.Be required. datacenter
* The third group，The length is not fixed.Be required. machine
* The fourth group，The length of the 38-42.Be required.The value of the current timestamp minus the fixed timestamp
* The fifth group，The length is not fixed.You don't need to fill it out. It's automatic。Sequence growth at the current time

## Installation
```shell
go get github.com/agclqq/snowflake
```
## usage
```go
package main

import (
	"fmt"
	"github.com/agclqq/snowflake"
)

func main()  {
	sf,err:=snowflake.New(2,2,2,2,snowflake.T38)
	if err!=nil{
		fmt.Println(err)
		return
	}
	id:=sf.GetId()
	fmt.Println(id)
}
```

## best practice

Using the TestSnowFlake_simg_GetId() method in the test case, you can measure the num value per second that you can achieve by adjusting the parameter `New()` and the `num variable` in the program to confirm your desired performance.


In this test case, the hardware environment of 'Intel(R) Core(TM) I7-10750H CPU @ 2.60GHz' can produce about 2.4 million serial numbers per second.


