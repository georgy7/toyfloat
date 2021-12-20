# Toyfloat

12-bit floating-point library.

Expected applications:

* File Format Design
* Storing colors (pixels, voxels)

It has:

* 7 bits normalized significand
* 4 bits exponent
* 1 sign bit
* (-255.9961, 255.9961) values range
* exact 0, 1, -1
* no infinities
* no NaN

![Formula](images/formula.png)

![Representation in memory](images/bits.png)
