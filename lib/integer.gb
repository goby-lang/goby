class Integer
  def times
    range = (0..self-1)

    if block_given?
      range.each do |i|
        yield(i)
      end
    end

    range.to_enum
  end
end
