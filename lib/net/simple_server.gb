module Net
  class SimpleServer
    attr_reader   :port
    attr_accessor :file_root

    def initialize(port)
      @port = port
    end

    def get(path)
      mount(path, "GET") do |req, res|
        yield(req, res)
      end
    end

    def post(path)
      mount(path, "POST") do |req, res|
        yield(req, res)
      end
    end

    def put(path)
      mount(path, "PUT") do |req, res|
        yield(req, res)
      end
    end

    def delete(path)
      mount(path, "DELETE") do |req, res|
        yield(req, res)
      end
    end

    def head(path)
      mount(path, "HEAD") do |req, res|
        yield(req, res)
      end
    end
  end
end
