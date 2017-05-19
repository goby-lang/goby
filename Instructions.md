## Overview

|Name                |arg1                  |arg2          |Example                   |Description|
|---                 |---                   |---           |---                       |---|
|getlocal            |depth                 |index         |`getlocal 0 1`            |Retrieve local variable|
|getconstant         |constant name         |              |`getconstant Bar`         |Retrive constant |
|getinstancevariable |instance variable name|              |`getinstancevariable @foo`|Retrieve instance variable|
|setlocal            |depth                 |index         |`setlocal 0 1`            |Set stack's top value to target location|
|setconstant         |constant name         |              |`setconstant Foo`         |Set stack's top value as constant|
|setinstancevariable |instance variable name|              |`setinstancevariable @bar`|Set stack's top value as instance variable|
|putstring           |string                |              |`putstring "Hello"`       |Put given string on the stack|
|putself             |                      |              |`putself`                 |Put self on the stack|
|putobject           |obj                   |              |`putobject 1`             |Put given object on the stack|
|putnil              |                      |              |`putnil`                  |Put nil on the stack|
|newarray            |object count          |              |`newarray 5`              |Initialize an array with last (object count) values on the stack|
|newhash             |key value count       |              |`newhash 10`              |Initialize a hash with last (key value count/2) pairs of key/values|
|branchunless        |location              |              |`branchunless 10`         |Jump to given location if stack's top value is false|
|branchif            |location              |              |`branchif 5`              |Jump to given location if stack's top value is true|
|jump                |location              |              |`jump 3`                  |Jump to given location|
|def_method          |parameters count      |              |`def_method 3`            |Define a method with given parameters count and uses stack's top value as method name|
|def_singleton_method|parameters count      |              |`def_singleton_method 2`  |Define a method with given parameters count and uses stack's top value as method name|
|def_class           |type:name             |superclass    |`def_class class:Foo Bar` |Define a class with given type, name and superclass(if it inherits from other class|
|send                |method name           |argument count|`send foo 2`              |Send a method with given name and arguments count|
|invokeblock         |argument count        |              |`invokeblock 1`           |Execute block with given arguments count|
|pop                 |                      |              |`pop`                     |Remove and return stack's top value and minus SP with 1|
|leave               |                      |              |`leave`                   |Finish current callframe's execution and pop it|