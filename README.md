# Benchmarks for XLSX library
It was not a goal to make best of the best, but the same time it's interesting to know pros/cons. 
For some cases this library is second, for other - best, but in case of reading huge files - **the only**. 

|                | tealeg | excelize | xlsx |
|----------------|:------:|:--------:|:----:|
| RandomGet      |   1!   |     3    |   2  |
| RandomSet      |   1!   |     3    |   2  |
| RandomSetStyle |   1!   |     3    |   2  |
| ReadBigFile    |    2   |     3    |   1  |
| UpdateBigFile  |    2!! |     3    |   1  |
| ReadHugeFile   |    -   |     -    |   1  |
| UpdateHugeFile |    -   |     -    |   1  |

* ! - does not mutate information directly, so faster get/set, but slower read/write files - sometimes it can take forever to open file.
* !! - corrupted file after saving, lost styles/formatting
 
[Benchmarks report](BENCHMARKS.md)