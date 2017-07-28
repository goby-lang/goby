## Steps to test Goby REPL(Igb)

Currently we are unable to perform automated tests against REPL. Follow the steps to test REPL manually.

* green promot: `» `
* red prompt: `¤ ` -- always indented
* yellow prompt: `#» ` -- result line

### 1. Basic operations

1. `goby -i` to start Igb.
    * expect:
        * startup massage with version no. and fortune are shown
        * fortune characters are random
        * green prompt `»` are shown
    ```ruby
    Goby 0.0.9 😽 😉 🤓
    »
    ```
2. type `help` and Return key
    * expect: following help messages are shown:
    ```ruby
    » help
    commands:
       help
       reset
       exit
    »
    ```
3. type `h` and then type Tab key
    * expect: autocomplete for `help` works
4. type `reset` and Return key
    * expect: following messages are shown
    ```ruby
    Restarting Igb...
    Goby 0.0.9 😎 😛 😪
    »
    ```
5. type `r` and then type Tab key
    * expect: autocomplet for `reset` works

6. perform shortcuts from [readline](https://github.com/chzyer/readline/blob/master/doc/shortcut.md)
    * expect: any shortcuts, including command history, are available

7. type `exit` and Return key
    * expect: the following message is shown and exited
    ```ruby
    » exit
    Bye!
    ```

### 2. trailing `;` feature for supressing echo back

1. type a statement with a trailing `;`
    * expect: echo back is suppressed:
    ```ruby
    » 5*7;
    »
    ```

2. type a block with a trailing `;`
    * expect: echo back is suppressed:
    ```ruby
    » def foo
    ¤   42
    » end;
    »
    ```

3. type a block with some `;` signs except at the end
    * expect: echo back is not suppressed
    ```ruby
    » def foo;
    ¤   42;
    » end
    #»
    »
    ```

4. type a sentence with arguments plus `;`, but no parens
    * expect: echo back is suppressed:
    ```ruby
    » puts 7*8;
    »
    ```

### 4. comments

1. type a leading comment
    * expect: just echoes back an empty yellow prompt:
    ```ruby
    » # test
    #»
    »
    ```

2. type a trailing comment
    * expect: just echoes back an empty yellow prompt:
    ```ruby
    » # test
    #»
    »
    ```

3. type a block with multiple types of comments
    * expect: works as follows, with the single `#` are correctly indented, and echoes back an empty yellow promot:
    ```ruby
    » def too #
    ¤   #
    ¤   42 #
    ¤   #
    » end #
    #»
    ```

 4. Paste " ¤ #"
    * expect: no error causes

### 5. Interruption

1. type Ctrl-z
    * expect: nothing happened (suppressed)
2. type Ctrl-c on green prompt, with no string entered
    * expect: works the same as `exit`
3. type Ctrl-c on green promot, with some strings entered
    * expect: "-- block cleared" are shown on the line, new green prompt remains
    in the following case, pressed ctrl just after `aa`:
    ```ruby
    » aa -- block cleared
    »
    ```
4. type Ctrl-c on red prompt, with some strings entered
    * expect: "-- block cleared" are shown on the line, red prompt remains and green prompt is shown in the next line. Then you can continue to enter, with the previous valid lines preserved.
    in the following case, pressed ctrl-c just after `aa`:
    ```ruby
    » def foo
    ¤   aa -- block cleared
    »
    ```
4. type Ctrl-c on red prompt, with no string entered
    * expect: "-- block cleared" are shown on the line, red prompt remains and green prompt is shown in the next line. Then you can continue to enter, with the previous valid lines preserved.
    in the following case, pressed ctrl-c on the empty red prompt line.
    ```ruby
    » a = 1
    #»
    » def foo
    ¤   -- block cleared
    » a
    #» 1
    »
    ```

5. type Ctrl-c on red prompt, in the midlle of the nested block
    * expect: "-- block cleared" are shown on the line, red prompt remains and green prompt is shown in the next line. Then you can continue to enter, with the previous valid lines preserved.
    in the following case, pressed ctrl-c on the empty red prompt line.
    ```ruby
    » a = 1
    #»
    » class Foo
    ¤   def bar
    ¤     puts "haha"
    ¤     if true
    ¤       puts "hehe"
    ¤     end
    ¤     puts "hoho"
    ¤      -- block cleared
    » a
    #» 1
    »
    ```

### 6. Pasting multiple lines

1. copy the following script and pasted to Igb:
    ```ruby
    x = 0
    y = 0

    while x < 10 do
      x = x + 1
      if x == 5
        next
      end
      y = y + 1
    end

    x + y
    ```

    * expect: obtains the following result:
    ```ruby
    » x = 0
    #»
    » y = 0
    #»
    »
    » while x < 10 do
    ¤   x = x + 1
    ¤   if x == 5
    ¤     next
    ¤   end
    ¤   y = y + 1
    » end
    #»
    »
    » x + y
    #» 19
    »
    ```

2. copy the following part of the strings from Igb and paste it to Igb again
    ```ruby
    » x = 0
    #»
    » y = 0
    #»
    »
    » while x < 10 do
    ¤   x = x + 1
    ¤   if x == 5
    ¤     next
    ¤   end
    ¤   y = y + 1
    » end
    #»
    »
    » x + y
    ```

    * expect: obtains the same result as 1., with all prompts truncated:
    ```ruby
    » x = 0
    #»
    » y = 0
    #»
    »
    » while x < 10 do
    ¤   x = x + 1
    ¤   if x == 5
    ¤     next
    ¤   end
    ¤   y = y + 1
    » end
    #»
    »
    » x + y
    #» 19
    »
    ```
