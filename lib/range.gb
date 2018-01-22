class Range
  def lazy
    LazyEnumerator.new(to_enum)
  end

  def to_enum
    RangeEnumerator.new(self)
  end
end
