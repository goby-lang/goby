class Array
  def include?(x)
    any? do |i|
      i == x
    end
  end

  def to_enum
    ArrayEnumerator.new(self)
  end

  # Return a lazy iterator for self.
  #
  def lazy
    LazyEnumerator.new(to_enum)
  end
end
