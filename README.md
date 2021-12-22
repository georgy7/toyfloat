# Toyfloat

12-bit floating-point library.

Expected applications:

* file format design
* lossy compression

It has:

* 7 bits normalized significand
* 4 bits exponent
* 1 sign bit
* (-255.9961, 255.9961) values range
* exact 0, 1, -1
* no NaN

![Formula](images/formula.png)

![Toyfloat in uint16: 4 empty bits, exponent, sign and mantissa.](images/bits.png)

![Precision graph](images/precision.png)

## Other options

Unsigned 12-bit format (utoyfloat)

`_ _ _ _ x x x x` `m m m m m m m m`

Signed 13-bit (toyfloat13)

`_ _ _ s x x x x` `m m m m m m m m`

Signed 14-bit (toyfloat14)

`_ _ x x x x s m` `m m m m m m m m`
