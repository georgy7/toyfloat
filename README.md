# Toyfloat

12-bit floating-point library.

Expected applications:

* File format design
* Lossy compression

It has:

* 7 bits normalized significand
* 4 bits exponent
* 1 sign bit
* (-255.9961, 255.9961) values range
* exact 0, 1, -1
* no NaN

![Formula](images/formula.png)

![Representation in memory](images/bits.png)

![Precision graph](images/precision.png)
