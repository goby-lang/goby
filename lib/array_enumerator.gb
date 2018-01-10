# Minimail implementation of an array enumerator.
#
# Assumes that the Enumerator interface has two methods: #has_next? and #next.
#
class ArrayEnumerator
  def initialize(array)
    @array = array
    @current_position = -1
  end

  # Returns true if there is another element is available.
  #
  def has_next?
    @current_position + 1 < @array.length
  end

  # Returns the next element, and advances the internal position.
  #
  # Raises an error if there are no elements available.
  #
  def next
    if !has_next?
      raise StopIteration, "No more elements!"
    end

    @current_position += 1

    @array[@current_position]
  end
end
