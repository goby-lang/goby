module Net
  class HTTP
    class Response
      attr_accessor :body, :status, :status_code, :protocol, :transfer_encoding, :http_version, :request_http_version, :request
      attr_reader :headers

      def initialize(headers = {})
        @headers = headers
      end

      def set_header(key, value)
        if @headers.nil?
          @headers = {}
        end
        @headers[key] = value
      end

      def get_header(key)
        @headers[key]
      end

      def remove_header(key)
        @headers.delete(key)
      end
    end
  end
end