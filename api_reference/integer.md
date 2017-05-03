# Integer

[Source code](https://github.com/rooby-lang/rooby/blob/master/vm/integer.go)

### + (Integer)

Returns the sum of self and another `Integer`.

### - (Integer)

Returns the subtraction of another `Integer` from self.

### * (Integer)

Returns self multiplying another `Integer`.

### ** (Integer)

Returns self squaring another `Integer`.

### / (Integer)

Returns self divided by another `Integer`.

### > (Integer)

Returns if self is larger than another `Integer`.

### >= (Integer)

Returns if self is larger than or equals to another `Integer`.

### < (Integer)

Returns if self is smaller than another `Integer`.

### <= (Integer)

Returns if self is smaller than or equals to another `Integer`.

### <=> (Integer)

Returns 1 if self is larger than the incoming `Integer`, -1 if smaller. Otherwise 0.

### == (Integer)

Returns if self is equal to another `Integer`.

### != (Integer)

Returns if self is not equal to another `Integer`.

### ++ (Integer)

Adds 1 to self and returns.

### -- (Integer)

Substracts 1 from self and returns.

### to_s

Returns a `String` representation of self.

### even

Returns if self is even.

### odd

Returns if self is odd.

### next

Returns self + 1.

### pred

Returns self - 1.

### times (&block)

Yields a block a number of times equals to self.
